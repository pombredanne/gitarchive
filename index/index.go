package index

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type Index struct {
	db *sql.DB

	insertFetchQ, insertDepQ *sql.Stmt
	selectQ, latestQ         *sql.Stmt

	insertBlacklistQ, selectBlacklistQ *sql.Stmt
}

func Open(dataSourceName string) (*Index, error) {
	db, err := sql.Open("mysql", dataSourceName+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	i := &Index{db: db}

	query := `CREATE TABLE IF NOT EXISTS Fetches (
		Name VARCHAR(255) NOT NULL, INDEX (Name), Parent VARCHAR(255),
		Timestamp DATETIME, Refs JSON,
		PackID BIGINT UNIQUE KEY AUTO_INCREMENT, PackRef VARCHAR(255))`
	if _, err = db.Exec(query); err != nil {
		return nil, errors.Wrap(err, "failed to create Fetches")
	}

	query = `CREATE TABLE IF NOT EXISTS PackDeps (ID BIGINT, INDEX (ID), Dep BIGINT)`
	if _, err = db.Exec(query); err != nil {
		return nil, errors.Wrap(err, "failed to create PackDeps")
	}

	query = `CREATE TABLE IF NOT EXISTS Blacklist (
		Name VARCHAR(255) NOT NULL, INDEX (Name), Whitelisted BOOLEAN NOT NULL DEFAULT 0)`
	if _, err = db.Exec(query); err != nil {
		return nil, errors.Wrap(err, "failed to create Blacklist")
	}

	prepStmts := []struct {
		name **sql.Stmt
		sql  string
	}{
		{
			&i.insertFetchQ,
			`INSERT INTO Fetches (Name, Parent, Timestamp, Refs, PackRef) VALUES (?, ?, ?, ?, ?)`,
		},
		{
			&i.insertDepQ,
			`INSERT INTO PackDeps (ID, Dep) VALUES (?, ?)`,
		},
		{
			&i.latestQ,
			`SELECT Timestamp FROM Fetches WHERE Name = ? ORDER BY Timestamp DESC LIMIT 1`,
		},
		{
			&i.selectQ,
			`SELECT Parent, Refs, PackID FROM Fetches WHERE Name = ? ORDER BY Timestamp DESC LIMIT 1`,
		},
		{
			&i.insertBlacklistQ,
			`INSERT INTO Blacklist (Name) VALUES (?)`,
		},
		{
			&i.selectBlacklistQ,
			`SELECT Whitelisted FROM Blacklist WHERE Name = ?`,
		},
	}

	for _, x := range prepStmts {
		stmt, err := db.Prepare(x.sql)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prepare '%s'", x.sql)
		}
		*x.name = stmt
	}

	return i, nil
}

func (i *Index) AddFetch(name, parent string, timestamp time.Time,
	refs map[string]string, packRef string, packDeps []string) error {
	r, err := json.Marshal(refs)
	if err != nil {
		return err
	}
	res, err := i.insertFetchQ.Exec(name, parent, timestamp, r, packRef)
	if err != nil {
		return err
	}
	packID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	for _, dep := range packDeps {
		_, err := i.insertDepQ.Exec(packID, dep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Index) GetLatest(name string) (timestamp time.Time, err error) {
	err = i.latestQ.QueryRow(name).Scan(&timestamp)
	if err == sql.ErrNoRows {
		timestamp = time.Time{}
		err = nil
	}
	return
}

func (i *Index) GetHaves(name string) (haves map[string]struct{}, deps []string, err error) {
	var parent, packID string
	var refs []byte
	err = i.selectQ.QueryRow(name).Scan(&parent, &refs, &packID)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		return
	}
	var r map[string]string
	err = json.Unmarshal(refs, &r)
	if err != nil {
		return
	}

	haves = make(map[string]struct{})
	for _, ref := range r {
		haves[ref] = struct{}{}
	}
	deps = append(deps, packID)

	if parent != "" {
		err = i.selectQ.QueryRow(parent).Scan(&parent, &refs, &packID)
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		if err != nil {
			return
		}
		err = json.Unmarshal(refs, &r)
		if err != nil {
			return
		}

		for _, ref := range r {
			haves[ref] = struct{}{}
		}
		deps = append(deps, packID)
	}

	return
}

func (i *Index) AddBlacklist(name string) error {
	_, err := i.insertBlacklistQ.Exec(name)
	return err
}

type BlacklistState int

const (
	Blacklisted BlacklistState = iota
	Whitelisted
	Neutral
)

func (i *Index) BlacklistState(name string) (BlacklistState, error) {
	var whitelisted bool
	err := i.selectBlacklistQ.QueryRow(name).Scan(&whitelisted)
	if err == sql.ErrNoRows {
		return Neutral, nil
	}
	if err != nil {
		return 0, err
	}
	if whitelisted {
		return Whitelisted, nil
	}
	return Blacklisted, nil
}

func (i *Index) Close() error {
	return i.db.Close()
}
