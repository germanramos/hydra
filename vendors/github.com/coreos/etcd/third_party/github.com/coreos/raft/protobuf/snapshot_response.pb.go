// Code generated by protoc-gen-gogo.
// source: snapshot_response.proto
// DO NOT EDIT!

package protobuf

import proto "github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// discarding unused import gogoproto "code.google.com/p/gogoprotobuf/gogoproto/gogo.pb"

import io8 "io"
import code_google_com_p_gogoprotobuf_proto16 "github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/code.google.com/p/gogoprotobuf/proto"

import fmt24 "fmt"
import strings16 "strings"
import reflect16 "reflect"

import fmt25 "fmt"
import strings17 "strings"
import code_google_com_p_gogoprotobuf_proto17 "github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/code.google.com/p/gogoprotobuf/proto"
import sort8 "sort"
import strconv8 "strconv"
import reflect17 "reflect"

import fmt26 "fmt"
import bytes8 "bytes"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type SnapshotResponse struct {
	Success			*bool	`protobuf:"varint,1,req" json:"Success,omitempty"`
	XXX_unrecognized	[]byte	`json:"-"`
}

func (m *SnapshotResponse) Reset()	{ *m = SnapshotResponse{} }
func (*SnapshotResponse) ProtoMessage()	{}

func (m *SnapshotResponse) GetSuccess() bool {
	if m != nil && m.Success != nil {
		return *m.Success
	}
	return false
}

func init() {
}
func (m *SnapshotResponse) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io8.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return proto.ErrWrongType
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io8.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.Success = &b
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := code_google_com_p_gogoprotobuf_proto16.Skip(data[index:])
			if err != nil {
				return err
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}
func (this *SnapshotResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings16.Join([]string{`&SnapshotResponse{`,
		`Success:` + valueToStringSnapshotResponse(this.Success) + `,`,
		`XXX_unrecognized:` + fmt24.Sprintf("%v", this.XXX_unrecognized) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringSnapshotResponse(v interface{}) string {
	rv := reflect16.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect16.Indirect(rv).Interface()
	return fmt24.Sprintf("*%v", pv)
}
func (m *SnapshotResponse) Size() (n int) {
	var l int
	_ = l
	if m.Success != nil {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovSnapshotResponse(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozSnapshotResponse(x uint64) (n int) {
	return sovSnapshotResponse(uint64((x << 1) ^ uint64((int64(x) >> 63))))
	return sovSnapshotResponse(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func NewPopulatedSnapshotResponse(r randySnapshotResponse, easy bool) *SnapshotResponse {
	this := &SnapshotResponse{}
	v1 := bool(r.Intn(2) == 0)
	this.Success = &v1
	if !easy && r.Intn(10) != 0 {
		this.XXX_unrecognized = randUnrecognizedSnapshotResponse(r, 2)
	}
	return this
}

type randySnapshotResponse interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneSnapshotResponse(r randySnapshotResponse) rune {
	res := rune(r.Uint32() % 1112064)
	if 55296 <= res {
		res += 2047
	}
	return res
}
func randStringSnapshotResponse(r randySnapshotResponse) string {
	v2 := r.Intn(100)
	tmps := make([]rune, v2)
	for i := 0; i < v2; i++ {
		tmps[i] = randUTF8RuneSnapshotResponse(r)
	}
	return string(tmps)
}
func randUnrecognizedSnapshotResponse(r randySnapshotResponse, maxFieldNumber int) (data []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		data = randFieldSnapshotResponse(data, r, fieldNumber, wire)
	}
	return data
}
func randFieldSnapshotResponse(data []byte, r randySnapshotResponse, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		data = encodeVarintPopulateSnapshotResponse(data, uint64(key))
		data = encodeVarintPopulateSnapshotResponse(data, uint64(r.Int63()))
	case 1:
		data = encodeVarintPopulateSnapshotResponse(data, uint64(key))
		data = append(data, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		data = encodeVarintPopulateSnapshotResponse(data, uint64(key))
		ll := r.Intn(100)
		data = encodeVarintPopulateSnapshotResponse(data, uint64(ll))
		for j := 0; j < ll; j++ {
			data = append(data, byte(r.Intn(256)))
		}
	default:
		data = encodeVarintPopulateSnapshotResponse(data, uint64(key))
		data = append(data, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return data
}
func encodeVarintPopulateSnapshotResponse(data []byte, v uint64) []byte {
	for v >= 1<<7 {
		data = append(data, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	data = append(data, uint8(v))
	return data
}
func (m *SnapshotResponse) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *SnapshotResponse) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Success != nil {
		data[i] = 0x8
		i++
		if *m.Success {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}
func encodeFixed64SnapshotResponse(data []byte, offset int, v uint64) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	data[offset+4] = uint8(v >> 32)
	data[offset+5] = uint8(v >> 40)
	data[offset+6] = uint8(v >> 48)
	data[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32SnapshotResponse(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintSnapshotResponse(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (this *SnapshotResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings17.Join([]string{`&protobuf.SnapshotResponse{` + `Success:` + valueToGoStringSnapshotResponse(this.Success, "bool"), `XXX_unrecognized:` + fmt25.Sprintf("%#v", this.XXX_unrecognized) + `}`}, ", ")
	return s
}
func valueToGoStringSnapshotResponse(v interface{}, typ string) string {
	rv := reflect17.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect17.Indirect(rv).Interface()
	return fmt25.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func extensionToGoStringSnapshotResponse(e map[int32]code_google_com_p_gogoprotobuf_proto17.Extension) string {
	if e == nil {
		return "nil"
	}
	s := "map[int32]proto.Extension{"
	keys := make([]int, 0, len(e))
	for k := range e {
		keys = append(keys, int(k))
	}
	sort8.Ints(keys)
	ss := []string{}
	for _, k := range keys {
		ss = append(ss, strconv8.Itoa(k)+": "+e[int32(k)].GoString())
	}
	s += strings17.Join(ss, ",") + "}"
	return s
}
func (this *SnapshotResponse) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt26.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*SnapshotResponse)
	if !ok {
		return fmt26.Errorf("that is not of type *SnapshotResponse")
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt26.Errorf("that is type *SnapshotResponse but is nil && this != nil")
	} else if this == nil {
		return fmt26.Errorf("that is type *SnapshotResponsebut is not nil && this == nil")
	}
	if this.Success != nil && that1.Success != nil {
		if *this.Success != *that1.Success {
			return fmt26.Errorf("Success this(%v) Not Equal that(%v)", *this.Success, *that1.Success)
		}
	} else if this.Success != nil {
		return fmt26.Errorf("this.Success == nil && that.Success != nil")
	} else if that1.Success != nil {
		return fmt26.Errorf("Success this(%v) Not Equal that(%v)", this.Success, that1.Success)
	}
	if !bytes8.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return fmt26.Errorf("XXX_unrecognized this(%v) Not Equal that(%v)", this.XXX_unrecognized, that1.XXX_unrecognized)
	}
	return nil
}
func (this *SnapshotResponse) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*SnapshotResponse)
	if !ok {
		return false
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Success != nil && that1.Success != nil {
		if *this.Success != *that1.Success {
			return false
		}
	} else if this.Success != nil {
		return false
	} else if that1.Success != nil {
		return false
	}
	if !bytes8.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
