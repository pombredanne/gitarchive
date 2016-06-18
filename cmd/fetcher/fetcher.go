package main

import (
	"expvar"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"

	"github.com/thecodearchive/gitarchive/git"
	"github.com/thecodearchive/gitarchive/index"
	"github.com/thecodearchive/gitarchive/queue"
	"github.com/thecodearchive/gitarchive/weekmap"
)

type Fetcher struct {
	q        *queue.Queue
	i        *index.Index
	bucket   *storage.BucketHandle
	schedule *weekmap.WeekMap

	exp *expvar.Map

	closing uint32
}

func (f *Fetcher) Run() error {
	f.exp.Set("fetchbytes", &expvar.Int{})
	for atomic.LoadUint32(&f.closing) == 0 {
		if !f.schedule.Get(time.Now()) {
			f.exp.Add("sleep", 1)
			interruptableSleep(5 * time.Minute)
			continue
		}

		name, parent, err := f.q.Pop()
		if err != nil {
			return err
		}

		if name == "" {
			f.exp.Add("emptyqueue", 1)
			interruptableSleep(30 * time.Second)
			continue
		}

		if err := f.Fetch(name, parent); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fetcher) Fetch(name, parent string) error {
	f.exp.Add("fetches", 1)

	url := "https://github.com/" + name + ".git"
	haves, deps, err := f.i.GetHaves(url)
	if err != nil {
		return err
	}

	if haves == nil {
		f.exp.Add("new", 1)
	}
	if parent != "" {
		f.exp.Add("forks", 1)
	}

	logVerb, logFork := "Cloning", ""
	if haves != nil {
		logVerb = "Fetching"
	}
	if parent != "" {
		logFork = fmt.Sprintf(" (fork of %s)", parent)
	}
	log.Printf("[+] %s %s%s...", logVerb, name, logFork)

	start := time.Now()
	bw := f.exp.Get("fetchbytes").(*expvar.Int)
	refs, r, err := git.Fetch(url, haves, os.Stderr, bw)
	if err != nil {
		return err
	}

	packRefName := fmt.Sprintf("%s|%d", url, time.Now().UnixNano())
	if r != nil {
		w := f.bucket.Object(packRefName).NewWriter(context.Background())

		bytesFetched, err := io.Copy(w, r)
		if err != nil {
			return err
		}
		w.Close()
		r.Close()
		f.exp.Add("fetchtime", int64(time.Since(start)))
		log.Printf("[+] Got %d refs, %d bytes in %s.", len(refs), bytesFetched, time.Since(start))
	} else {
		// Empty packfile.
		packRefName = "EMPTY|" + packRefName
		f.exp.Add("emptypack", 1)
		log.Printf("[+] Got %d refs, and a empty packfile.", len(refs))
	}

	urlParent := ""
	if parent != "" {
		urlParent = "https://github.com/" + parent + ".git"
	}

	f.i.AddFetch(url, urlParent, time.Now(), refs, packRefName, deps)
	return err
}

func (f *Fetcher) Stop() {
	atomic.StoreUint32(&f.closing, 1)
}

func interruptableSleep(d time.Duration) bool {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(c)
	t := time.NewTimer(d)
	select {
	case <-c:
		return false
	case <-t.C:
		return true
	}
}
