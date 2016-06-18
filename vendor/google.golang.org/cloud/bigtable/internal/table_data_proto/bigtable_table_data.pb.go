// Code generated by protoc-gen-go.
// source: google.golang.org/cloud/bigtable/internal/table_data_proto/bigtable_table_data.proto
// DO NOT EDIT!

/*
Package google_bigtable_admin_table_v1 is a generated protocol buffer package.

It is generated from these files:
	google.golang.org/cloud/bigtable/internal/table_data_proto/bigtable_table_data.proto

It has these top-level messages:
	Table
	ColumnFamily
	GcRule
*/
package google_bigtable_admin_table_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/duration"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Table_TimestampGranularity int32

const (
	Table_MILLIS Table_TimestampGranularity = 0
)

var Table_TimestampGranularity_name = map[int32]string{
	0: "MILLIS",
}
var Table_TimestampGranularity_value = map[string]int32{
	"MILLIS": 0,
}

func (x Table_TimestampGranularity) String() string {
	return proto.EnumName(Table_TimestampGranularity_name, int32(x))
}
func (Table_TimestampGranularity) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 0}
}

// A collection of user data indexed by row, column, and timestamp.
// Each table is served using the resources of its parent cluster.
type Table struct {
	// A unique identifier of the form
	// <cluster_name>/tables/[_a-zA-Z0-9][-_.a-zA-Z0-9]*
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// The column families configured for this table, mapped by column family id.
	ColumnFamilies map[string]*ColumnFamily `protobuf:"bytes,3,rep,name=column_families,json=columnFamilies" json:"column_families,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// The granularity (e.g. MILLIS, MICROS) at which timestamps are stored in
	// this table. Timestamps not matching the granularity will be rejected.
	// Cannot be changed once the table is created.
	Granularity Table_TimestampGranularity `protobuf:"varint,4,opt,name=granularity,enum=google.bigtable.admin.table.v1.Table_TimestampGranularity" json:"granularity,omitempty"`
}

func (m *Table) Reset()                    { *m = Table{} }
func (m *Table) String() string            { return proto.CompactTextString(m) }
func (*Table) ProtoMessage()               {}
func (*Table) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Table) GetColumnFamilies() map[string]*ColumnFamily {
	if m != nil {
		return m.ColumnFamilies
	}
	return nil
}

// A set of columns within a table which share a common configuration.
type ColumnFamily struct {
	// A unique identifier of the form <table_name>/columnFamilies/[-_.a-zA-Z0-9]+
	// The last segment is the same as the "name" field in
	// google.bigtable.v1.Family.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Garbage collection expression specified by the following grammar:
	//   GC = EXPR
	//      | "" ;
	//   EXPR = EXPR, "||", EXPR              (* lowest precedence *)
	//        | EXPR, "&&", EXPR
	//        | "(", EXPR, ")"                (* highest precedence *)
	//        | PROP ;
	//   PROP = "version() >", NUM32
	//        | "age() >", NUM64, [ UNIT ] ;
	//   NUM32 = non-zero-digit { digit } ;    (* # NUM32 <= 2^32 - 1 *)
	//   NUM64 = non-zero-digit { digit } ;    (* # NUM64 <= 2^63 - 1 *)
	//   UNIT =  "d" | "h" | "m"  (* d=days, h=hours, m=minutes, else micros *)
	// GC expressions can be up to 500 characters in length
	//
	// The different types of PROP are defined as follows:
	//   version() - cell index, counting from most recent and starting at 1
	//   age() - age of the cell (current time minus cell timestamp)
	//
	// Example: "version() > 3 || (age() > 3d && version() > 1)"
	//   drop cells beyond the most recent three, and drop cells older than three
	//   days unless they're the most recent cell in the row/column
	//
	// Garbage collection executes opportunistically in the background, and so
	// it's possible for reads to return a cell even if it matches the active GC
	// expression for its family.
	GcExpression string `protobuf:"bytes,2,opt,name=gc_expression,json=gcExpression" json:"gc_expression,omitempty"`
	// Garbage collection rule specified as a protobuf.
	// Supersedes `gc_expression`.
	// Must serialize to at most 500 bytes.
	//
	// NOTE: Garbage collection executes opportunistically in the background, and
	// so it's possible for reads to return a cell even if it matches the active
	// GC expression for its family.
	GcRule *GcRule `protobuf:"bytes,3,opt,name=gc_rule,json=gcRule" json:"gc_rule,omitempty"`
}

func (m *ColumnFamily) Reset()                    { *m = ColumnFamily{} }
func (m *ColumnFamily) String() string            { return proto.CompactTextString(m) }
func (*ColumnFamily) ProtoMessage()               {}
func (*ColumnFamily) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ColumnFamily) GetGcRule() *GcRule {
	if m != nil {
		return m.GcRule
	}
	return nil
}

// Rule for determining which cells to delete during garbage collection.
type GcRule struct {
	// Types that are valid to be assigned to Rule:
	//	*GcRule_MaxNumVersions
	//	*GcRule_MaxAge
	//	*GcRule_Intersection_
	//	*GcRule_Union_
	Rule isGcRule_Rule `protobuf_oneof:"rule"`
}

func (m *GcRule) Reset()                    { *m = GcRule{} }
func (m *GcRule) String() string            { return proto.CompactTextString(m) }
func (*GcRule) ProtoMessage()               {}
func (*GcRule) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type isGcRule_Rule interface {
	isGcRule_Rule()
}

type GcRule_MaxNumVersions struct {
	MaxNumVersions int32 `protobuf:"varint,1,opt,name=max_num_versions,json=maxNumVersions,oneof"`
}
type GcRule_MaxAge struct {
	MaxAge *google_protobuf.Duration `protobuf:"bytes,2,opt,name=max_age,json=maxAge,oneof"`
}
type GcRule_Intersection_ struct {
	Intersection *GcRule_Intersection `protobuf:"bytes,3,opt,name=intersection,oneof"`
}
type GcRule_Union_ struct {
	Union *GcRule_Union `protobuf:"bytes,4,opt,name=union,oneof"`
}

func (*GcRule_MaxNumVersions) isGcRule_Rule() {}
func (*GcRule_MaxAge) isGcRule_Rule()         {}
func (*GcRule_Intersection_) isGcRule_Rule()  {}
func (*GcRule_Union_) isGcRule_Rule()         {}

func (m *GcRule) GetRule() isGcRule_Rule {
	if m != nil {
		return m.Rule
	}
	return nil
}

func (m *GcRule) GetMaxNumVersions() int32 {
	if x, ok := m.GetRule().(*GcRule_MaxNumVersions); ok {
		return x.MaxNumVersions
	}
	return 0
}

func (m *GcRule) GetMaxAge() *google_protobuf.Duration {
	if x, ok := m.GetRule().(*GcRule_MaxAge); ok {
		return x.MaxAge
	}
	return nil
}

func (m *GcRule) GetIntersection() *GcRule_Intersection {
	if x, ok := m.GetRule().(*GcRule_Intersection_); ok {
		return x.Intersection
	}
	return nil
}

func (m *GcRule) GetUnion() *GcRule_Union {
	if x, ok := m.GetRule().(*GcRule_Union_); ok {
		return x.Union
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*GcRule) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _GcRule_OneofMarshaler, _GcRule_OneofUnmarshaler, _GcRule_OneofSizer, []interface{}{
		(*GcRule_MaxNumVersions)(nil),
		(*GcRule_MaxAge)(nil),
		(*GcRule_Intersection_)(nil),
		(*GcRule_Union_)(nil),
	}
}

func _GcRule_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*GcRule)
	// rule
	switch x := m.Rule.(type) {
	case *GcRule_MaxNumVersions:
		b.EncodeVarint(1<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.MaxNumVersions))
	case *GcRule_MaxAge:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MaxAge); err != nil {
			return err
		}
	case *GcRule_Intersection_:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Intersection); err != nil {
			return err
		}
	case *GcRule_Union_:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Union); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("GcRule.Rule has unexpected type %T", x)
	}
	return nil
}

func _GcRule_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*GcRule)
	switch tag {
	case 1: // rule.max_num_versions
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Rule = &GcRule_MaxNumVersions{int32(x)}
		return true, err
	case 2: // rule.max_age
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf.Duration)
		err := b.DecodeMessage(msg)
		m.Rule = &GcRule_MaxAge{msg}
		return true, err
	case 3: // rule.intersection
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(GcRule_Intersection)
		err := b.DecodeMessage(msg)
		m.Rule = &GcRule_Intersection_{msg}
		return true, err
	case 4: // rule.union
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(GcRule_Union)
		err := b.DecodeMessage(msg)
		m.Rule = &GcRule_Union_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _GcRule_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*GcRule)
	// rule
	switch x := m.Rule.(type) {
	case *GcRule_MaxNumVersions:
		n += proto.SizeVarint(1<<3 | proto.WireVarint)
		n += proto.SizeVarint(uint64(x.MaxNumVersions))
	case *GcRule_MaxAge:
		s := proto.Size(x.MaxAge)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GcRule_Intersection_:
		s := proto.Size(x.Intersection)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GcRule_Union_:
		s := proto.Size(x.Union)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// A GcRule which deletes cells matching all of the given rules.
type GcRule_Intersection struct {
	// Only delete cells which would be deleted by every element of `rules`.
	Rules []*GcRule `protobuf:"bytes,1,rep,name=rules" json:"rules,omitempty"`
}

func (m *GcRule_Intersection) Reset()                    { *m = GcRule_Intersection{} }
func (m *GcRule_Intersection) String() string            { return proto.CompactTextString(m) }
func (*GcRule_Intersection) ProtoMessage()               {}
func (*GcRule_Intersection) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

func (m *GcRule_Intersection) GetRules() []*GcRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

// A GcRule which deletes cells matching any of the given rules.
type GcRule_Union struct {
	// Delete cells which would be deleted by any element of `rules`.
	Rules []*GcRule `protobuf:"bytes,1,rep,name=rules" json:"rules,omitempty"`
}

func (m *GcRule_Union) Reset()                    { *m = GcRule_Union{} }
func (m *GcRule_Union) String() string            { return proto.CompactTextString(m) }
func (*GcRule_Union) ProtoMessage()               {}
func (*GcRule_Union) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 1} }

func (m *GcRule_Union) GetRules() []*GcRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

func init() {
	proto.RegisterType((*Table)(nil), "google.bigtable.admin.table.v1.Table")
	proto.RegisterType((*ColumnFamily)(nil), "google.bigtable.admin.table.v1.ColumnFamily")
	proto.RegisterType((*GcRule)(nil), "google.bigtable.admin.table.v1.GcRule")
	proto.RegisterType((*GcRule_Intersection)(nil), "google.bigtable.admin.table.v1.GcRule.Intersection")
	proto.RegisterType((*GcRule_Union)(nil), "google.bigtable.admin.table.v1.GcRule.Union")
	proto.RegisterEnum("google.bigtable.admin.table.v1.Table_TimestampGranularity", Table_TimestampGranularity_name, Table_TimestampGranularity_value)
}

func init() {
	proto.RegisterFile("google.golang.org/cloud/bigtable/internal/table_data_proto/bigtable_table_data.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 519 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x54, 0x6f, 0x8b, 0xd3, 0x4e,
	0x10, 0xbe, 0xfe, 0x49, 0x8e, 0x9b, 0xf6, 0xd7, 0x5f, 0x59, 0x45, 0x6a, 0x5f, 0xc8, 0x11, 0x41,
	0x0e, 0x91, 0x04, 0x7b, 0xbe, 0xd0, 0x43, 0x14, 0x6b, 0xeb, 0x59, 0xa8, 0x72, 0xc4, 0x2a, 0x08,
	0x42, 0xd8, 0xa6, 0x7b, 0x4b, 0x30, 0xbb, 0x5b, 0x36, 0xd9, 0x72, 0xfd, 0x06, 0x7e, 0x14, 0x3f,
	0x9f, 0x9f, 0xc0, 0xdd, 0x4d, 0x7a, 0x97, 0x83, 0x62, 0x2b, 0xbe, 0xea, 0x74, 0xe6, 0x79, 0x9e,
	0x79, 0x32, 0x33, 0x09, 0xcc, 0xa8, 0x10, 0x34, 0x25, 0x3e, 0x15, 0x29, 0xe6, 0xd4, 0x17, 0x92,
	0x06, 0x71, 0x2a, 0xd4, 0x22, 0x98, 0x27, 0x34, 0xc7, 0xf3, 0x94, 0x04, 0x09, 0xcf, 0x89, 0xe4,
	0x38, 0x0d, 0xec, 0xdf, 0x68, 0x81, 0x73, 0x1c, 0x2d, 0xa5, 0xc8, 0xc5, 0x35, 0x24, 0xba, 0xa9,
	0xf8, 0xb6, 0x82, 0x1e, 0x94, 0xaa, 0x1b, 0x84, 0x8f, 0x17, 0x2c, 0xe1, 0x7e, 0x11, 0xaf, 0x9e,
	0xf6, 0xcb, 0x7a, 0x60, 0xd1, 0x73, 0x75, 0x19, 0x2c, 0x94, 0xc4, 0x79, 0x22, 0x78, 0xc1, 0xf7,
	0x7e, 0xd5, 0xc1, 0x99, 0x19, 0x30, 0x42, 0xd0, 0xe4, 0x98, 0x91, 0x5e, 0xed, 0xb8, 0x76, 0x72,
	0x14, 0xda, 0x18, 0xcd, 0xe1, 0xff, 0x58, 0xa4, 0x8a, 0xf1, 0xe8, 0x12, 0xb3, 0x24, 0x4d, 0x48,
	0xd6, 0x6b, 0x1c, 0x37, 0x4e, 0x5a, 0x83, 0x17, 0xfe, 0x9f, 0xfb, 0xfa, 0x56, 0xd3, 0x7f, 0x6b,
	0xc9, 0xef, 0x4a, 0xee, 0x98, 0xe7, 0x72, 0x1d, 0x76, 0xe2, 0x5b, 0x49, 0xf4, 0x0d, 0x5a, 0x54,
	0x62, 0xae, 0x52, 0x2c, 0x93, 0x7c, 0xdd, 0x6b, 0xea, 0xf6, 0x9d, 0xc1, 0xd9, 0x7e, 0xfa, 0xb3,
	0x84, 0x91, 0x2c, 0xc7, 0x6c, 0x79, 0x7e, 0xa3, 0x10, 0x56, 0xe5, 0xfa, 0x02, 0xee, 0x6c, 0x31,
	0x81, 0xba, 0xd0, 0xf8, 0x4e, 0xd6, 0xe5, 0xb3, 0x9a, 0x10, 0x0d, 0xc1, 0x59, 0xe1, 0x54, 0x91,
	0x5e, 0x5d, 0xe7, 0x5a, 0x83, 0x27, 0xbb, 0x0c, 0x54, 0x54, 0xd7, 0x61, 0x41, 0x3d, 0xab, 0x3f,
	0xaf, 0x79, 0x1e, 0xdc, 0xdd, 0xe6, 0x0a, 0x01, 0xb8, 0x1f, 0x26, 0xd3, 0xe9, 0xe4, 0x53, 0xf7,
	0xc0, 0xfb, 0x51, 0x83, 0x76, 0x95, 0xbf, 0x75, 0xf6, 0x0f, 0xe1, 0x3f, 0x1a, 0x47, 0xe4, 0x6a,
	0x29, 0x49, 0x96, 0xe9, 0x85, 0x59, 0x63, 0x47, 0x61, 0x9b, 0xc6, 0xe3, 0xeb, 0x1c, 0x7a, 0x0d,
	0x87, 0x1a, 0x24, 0x55, 0x4a, 0xf4, 0x62, 0x8c, 0xef, 0x47, 0xbb, 0x7c, 0x9f, 0xc7, 0xa1, 0x46,
	0x87, 0x2e, 0xb5, 0xbf, 0xde, 0xcf, 0x06, 0xb8, 0x45, 0x0a, 0x3d, 0x86, 0x2e, 0xc3, 0x57, 0x11,
	0x57, 0x2c, 0x5a, 0x11, 0x69, 0xe4, 0x33, 0x6b, 0xc8, 0x79, 0x7f, 0x10, 0x76, 0x74, 0xe5, 0xa3,
	0x62, 0x5f, 0xca, 0x3c, 0x7a, 0x06, 0x87, 0x06, 0x8b, 0xe9, 0x66, 0x5e, 0xf7, 0x37, 0x7d, 0x37,
	0x87, 0xe6, 0x8f, 0xca, 0x43, 0xd3, 0x6c, 0x57, 0x63, 0xdf, 0x50, 0x82, 0xbe, 0x42, 0xdb, 0xde,
	0x78, 0x46, 0x62, 0x53, 0x29, 0x2d, 0x9f, 0xee, 0x67, 0xd9, 0x9f, 0x54, 0xa8, 0x5a, 0xf4, 0x96,
	0x14, 0x1a, 0x81, 0xa3, 0xb8, 0xd1, 0x6c, 0xee, 0xb7, 0xbe, 0x52, 0xf3, 0x33, 0x2f, 0xc4, 0x0a,
	0x72, 0x7f, 0x0a, 0xed, 0x6a, 0x17, 0xf4, 0x12, 0x1c, 0x33, 0x5b, 0x33, 0x87, 0xc6, 0x5f, 0x0c,
	0xb7, 0x20, 0xf5, 0xc7, 0xe0, 0x58, 0xfd, 0x7f, 0x93, 0x19, 0xba, 0xd0, 0x34, 0xc1, 0xf0, 0x15,
	0x78, 0xb1, 0x60, 0x3b, 0xb8, 0xc3, 0x7b, 0xc3, 0xb2, 0x60, 0xdf, 0x90, 0x91, 0xfe, 0x52, 0x5c,
	0x98, 0x8d, 0x5c, 0xd4, 0xe6, 0xae, 0x5d, 0xcd, 0xe9, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa9,
	0x9e, 0x64, 0x6f, 0x89, 0x04, 0x00, 0x00,
}