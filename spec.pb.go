// Code generated by protoc-gen-gogo.
// source: spec.proto
// DO NOT EDIT!

/*
	Package main is a generated protocol buffer package.

	It is generated from these files:
		spec.proto

	It has these top-level messages:
		Events
		Event
		Field
		FileMonitor
		Meta
*/
package main

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import bytes "bytes"

import strings "strings"
import github_com_gogo_protobuf_proto "github.com/gogo/protobuf/proto"
import sort "sort"
import strconv "strconv"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Events struct {
	Bloom  []byte   `protobuf:"bytes,1,opt,name=bloom,proto3" json:"bloom,omitempty"`
	Events []*Event `protobuf:"bytes,2,rep,name=events" json:"events,omitempty"`
}

func (m *Events) Reset()                    { *m = Events{} }
func (*Events) ProtoMessage()               {}
func (*Events) Descriptor() ([]byte, []int) { return fileDescriptorSpec, []int{0} }

func (m *Events) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

type Event struct {
	Bloom  []byte   `protobuf:"bytes,1,opt,name=bloom,proto3" json:"bloom,omitempty"`
	Data   string   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Lines  int32    `protobuf:"varint,3,opt,name=lines,proto3" json:"lines,omitempty"`
	Path   string   `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	Fields []*Field `protobuf:"bytes,5,rep,name=fields" json:"fields,omitempty"`
	Ts     int64    `protobuf:"varint,6,opt,name=ts,proto3" json:"ts,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptorSpec, []int{1} }

func (m *Event) GetFields() []*Field {
	if m != nil {
		return m.Fields
	}
	return nil
}

type Field struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *Field) Reset()                    { *m = Field{} }
func (*Field) ProtoMessage()               {}
func (*Field) Descriptor() ([]byte, []int) { return fileDescriptorSpec, []int{2} }

type FileMonitor struct {
	Path   string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Offset int64  `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Poll   bool   `protobuf:"varint,3,opt,name=poll,proto3" json:"poll,omitempty"`
}

func (m *FileMonitor) Reset()                    { *m = FileMonitor{} }
func (*FileMonitor) ProtoMessage()               {}
func (*FileMonitor) Descriptor() ([]byte, []int) { return fileDescriptorSpec, []int{3} }

type Meta struct {
	Count int64 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (m *Meta) Reset()                    { *m = Meta{} }
func (*Meta) ProtoMessage()               {}
func (*Meta) Descriptor() ([]byte, []int) { return fileDescriptorSpec, []int{4} }

func init() {
	proto.RegisterType((*Events)(nil), "main.Events")
	proto.RegisterType((*Event)(nil), "main.Event")
	proto.RegisterType((*Field)(nil), "main.Field")
	proto.RegisterType((*FileMonitor)(nil), "main.FileMonitor")
	proto.RegisterType((*Meta)(nil), "main.Meta")
}
func (this *Events) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*Events)
	if !ok {
		that2, ok := that.(Events)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.Bloom, that1.Bloom) {
		return false
	}
	if len(this.Events) != len(that1.Events) {
		return false
	}
	for i := range this.Events {
		if !this.Events[i].Equal(that1.Events[i]) {
			return false
		}
	}
	return true
}
func (this *Event) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*Event)
	if !ok {
		that2, ok := that.(Event)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.Bloom, that1.Bloom) {
		return false
	}
	if this.Data != that1.Data {
		return false
	}
	if this.Lines != that1.Lines {
		return false
	}
	if this.Path != that1.Path {
		return false
	}
	if len(this.Fields) != len(that1.Fields) {
		return false
	}
	for i := range this.Fields {
		if !this.Fields[i].Equal(that1.Fields[i]) {
			return false
		}
	}
	if this.Ts != that1.Ts {
		return false
	}
	return true
}
func (this *Field) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*Field)
	if !ok {
		that2, ok := that.(Field)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Key != that1.Key {
		return false
	}
	if this.Value != that1.Value {
		return false
	}
	return true
}
func (this *FileMonitor) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*FileMonitor)
	if !ok {
		that2, ok := that.(FileMonitor)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Path != that1.Path {
		return false
	}
	if this.Offset != that1.Offset {
		return false
	}
	if this.Poll != that1.Poll {
		return false
	}
	return true
}
func (this *Meta) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*Meta)
	if !ok {
		that2, ok := that.(Meta)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Count != that1.Count {
		return false
	}
	return true
}
func (this *Events) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&main.Events{")
	s = append(s, "Bloom: "+fmt.Sprintf("%#v", this.Bloom)+",\n")
	if this.Events != nil {
		s = append(s, "Events: "+fmt.Sprintf("%#v", this.Events)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Event) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&main.Event{")
	s = append(s, "Bloom: "+fmt.Sprintf("%#v", this.Bloom)+",\n")
	s = append(s, "Data: "+fmt.Sprintf("%#v", this.Data)+",\n")
	s = append(s, "Lines: "+fmt.Sprintf("%#v", this.Lines)+",\n")
	s = append(s, "Path: "+fmt.Sprintf("%#v", this.Path)+",\n")
	if this.Fields != nil {
		s = append(s, "Fields: "+fmt.Sprintf("%#v", this.Fields)+",\n")
	}
	s = append(s, "Ts: "+fmt.Sprintf("%#v", this.Ts)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Field) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&main.Field{")
	s = append(s, "Key: "+fmt.Sprintf("%#v", this.Key)+",\n")
	s = append(s, "Value: "+fmt.Sprintf("%#v", this.Value)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *FileMonitor) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&main.FileMonitor{")
	s = append(s, "Path: "+fmt.Sprintf("%#v", this.Path)+",\n")
	s = append(s, "Offset: "+fmt.Sprintf("%#v", this.Offset)+",\n")
	s = append(s, "Poll: "+fmt.Sprintf("%#v", this.Poll)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Meta) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&main.Meta{")
	s = append(s, "Count: "+fmt.Sprintf("%#v", this.Count)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringSpec(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func extensionToGoStringSpec(m github_com_gogo_protobuf_proto.Message) string {
	e := github_com_gogo_protobuf_proto.GetUnsafeExtensionsMap(m)
	if e == nil {
		return "nil"
	}
	s := "proto.NewUnsafeXXX_InternalExtensions(map[int32]proto.Extension{"
	keys := make([]int, 0, len(e))
	for k := range e {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	ss := []string{}
	for _, k := range keys {
		ss = append(ss, strconv.Itoa(k)+": "+e[int32(k)].GoString())
	}
	s += strings.Join(ss, ",") + "})"
	return s
}
func (m *Events) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Events) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Bloom) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintSpec(data, i, uint64(len(m.Bloom)))
		i += copy(data[i:], m.Bloom)
	}
	if len(m.Events) > 0 {
		for _, msg := range m.Events {
			data[i] = 0x12
			i++
			i = encodeVarintSpec(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Event) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Event) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Bloom) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintSpec(data, i, uint64(len(m.Bloom)))
		i += copy(data[i:], m.Bloom)
	}
	if len(m.Data) > 0 {
		data[i] = 0x12
		i++
		i = encodeVarintSpec(data, i, uint64(len(m.Data)))
		i += copy(data[i:], m.Data)
	}
	if m.Lines != 0 {
		data[i] = 0x18
		i++
		i = encodeVarintSpec(data, i, uint64(m.Lines))
	}
	if len(m.Path) > 0 {
		data[i] = 0x22
		i++
		i = encodeVarintSpec(data, i, uint64(len(m.Path)))
		i += copy(data[i:], m.Path)
	}
	if len(m.Fields) > 0 {
		for _, msg := range m.Fields {
			data[i] = 0x2a
			i++
			i = encodeVarintSpec(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Ts != 0 {
		data[i] = 0x30
		i++
		i = encodeVarintSpec(data, i, uint64(m.Ts))
	}
	return i, nil
}

func (m *Field) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Field) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Key) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintSpec(data, i, uint64(len(m.Key)))
		i += copy(data[i:], m.Key)
	}
	if len(m.Value) > 0 {
		data[i] = 0x12
		i++
		i = encodeVarintSpec(data, i, uint64(len(m.Value)))
		i += copy(data[i:], m.Value)
	}
	return i, nil
}

func (m *FileMonitor) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *FileMonitor) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Path) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintSpec(data, i, uint64(len(m.Path)))
		i += copy(data[i:], m.Path)
	}
	if m.Offset != 0 {
		data[i] = 0x10
		i++
		i = encodeVarintSpec(data, i, uint64(m.Offset))
	}
	if m.Poll {
		data[i] = 0x18
		i++
		if m.Poll {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *Meta) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Meta) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Count != 0 {
		data[i] = 0x8
		i++
		i = encodeVarintSpec(data, i, uint64(m.Count))
	}
	return i, nil
}

func encodeFixed64Spec(data []byte, offset int, v uint64) int {
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
func encodeFixed32Spec(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintSpec(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (m *Events) Size() (n int) {
	var l int
	_ = l
	l = len(m.Bloom)
	if l > 0 {
		n += 1 + l + sovSpec(uint64(l))
	}
	if len(m.Events) > 0 {
		for _, e := range m.Events {
			l = e.Size()
			n += 1 + l + sovSpec(uint64(l))
		}
	}
	return n
}

func (m *Event) Size() (n int) {
	var l int
	_ = l
	l = len(m.Bloom)
	if l > 0 {
		n += 1 + l + sovSpec(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovSpec(uint64(l))
	}
	if m.Lines != 0 {
		n += 1 + sovSpec(uint64(m.Lines))
	}
	l = len(m.Path)
	if l > 0 {
		n += 1 + l + sovSpec(uint64(l))
	}
	if len(m.Fields) > 0 {
		for _, e := range m.Fields {
			l = e.Size()
			n += 1 + l + sovSpec(uint64(l))
		}
	}
	if m.Ts != 0 {
		n += 1 + sovSpec(uint64(m.Ts))
	}
	return n
}

func (m *Field) Size() (n int) {
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovSpec(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovSpec(uint64(l))
	}
	return n
}

func (m *FileMonitor) Size() (n int) {
	var l int
	_ = l
	l = len(m.Path)
	if l > 0 {
		n += 1 + l + sovSpec(uint64(l))
	}
	if m.Offset != 0 {
		n += 1 + sovSpec(uint64(m.Offset))
	}
	if m.Poll {
		n += 2
	}
	return n
}

func (m *Meta) Size() (n int) {
	var l int
	_ = l
	if m.Count != 0 {
		n += 1 + sovSpec(uint64(m.Count))
	}
	return n
}

func sovSpec(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozSpec(x uint64) (n int) {
	return sovSpec(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Events) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Events{`,
		`Bloom:` + fmt.Sprintf("%v", this.Bloom) + `,`,
		`Events:` + strings.Replace(fmt.Sprintf("%v", this.Events), "Event", "Event", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Event) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Event{`,
		`Bloom:` + fmt.Sprintf("%v", this.Bloom) + `,`,
		`Data:` + fmt.Sprintf("%v", this.Data) + `,`,
		`Lines:` + fmt.Sprintf("%v", this.Lines) + `,`,
		`Path:` + fmt.Sprintf("%v", this.Path) + `,`,
		`Fields:` + strings.Replace(fmt.Sprintf("%v", this.Fields), "Field", "Field", 1) + `,`,
		`Ts:` + fmt.Sprintf("%v", this.Ts) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Field) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Field{`,
		`Key:` + fmt.Sprintf("%v", this.Key) + `,`,
		`Value:` + fmt.Sprintf("%v", this.Value) + `,`,
		`}`,
	}, "")
	return s
}
func (this *FileMonitor) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&FileMonitor{`,
		`Path:` + fmt.Sprintf("%v", this.Path) + `,`,
		`Offset:` + fmt.Sprintf("%v", this.Offset) + `,`,
		`Poll:` + fmt.Sprintf("%v", this.Poll) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Meta) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Meta{`,
		`Count:` + fmt.Sprintf("%v", this.Count) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringSpec(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Events) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSpec
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Events: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Events: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bloom", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSpec
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bloom = append(m.Bloom[:0], data[iNdEx:postIndex]...)
			if m.Bloom == nil {
				m.Bloom = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Events", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSpec
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Events = append(m.Events, &Event{})
			if err := m.Events[len(m.Events)-1].Unmarshal(data[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSpec(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSpec
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Event) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSpec
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Event: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Event: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bloom", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSpec
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bloom = append(m.Bloom[:0], data[iNdEx:postIndex]...)
			if m.Bloom == nil {
				m.Bloom = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpec
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Lines", wireType)
			}
			m.Lines = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.Lines |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Path", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpec
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Path = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Fields", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSpec
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Fields = append(m.Fields, &Field{})
			if err := m.Fields[len(m.Fields)-1].Unmarshal(data[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ts", wireType)
			}
			m.Ts = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.Ts |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSpec(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSpec
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Field) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSpec
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Field: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Field: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpec
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpec
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSpec(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSpec
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FileMonitor) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSpec
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FileMonitor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FileMonitor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Path", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpec
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Path = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Offset", wireType)
			}
			m.Offset = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.Offset |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Poll", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Poll = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipSpec(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSpec
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Meta) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSpec
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Meta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Meta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			m.Count = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.Count |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSpec(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSpec
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSpec(data []byte) (n int, err error) {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSpec
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if data[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSpec
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthSpec
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowSpec
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := data[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipSpec(data[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthSpec = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSpec   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("spec.proto", fileDescriptorSpec) }

var fileDescriptorSpec = []byte{
	// 309 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x51, 0x31, 0x4e, 0xeb, 0x40,
	0x10, 0xcd, 0x66, 0x6d, 0xeb, 0x67, 0xf2, 0x85, 0xd0, 0x0a, 0x21, 0x17, 0xc8, 0x8a, 0x4c, 0xe3,
	0x02, 0x19, 0x09, 0x6e, 0x00, 0x22, 0x9d, 0x9b, 0xbd, 0x81, 0x93, 0xac, 0x85, 0xc5, 0xc6, 0x6b,
	0xc5, 0x9b, 0x48, 0x74, 0x1c, 0x21, 0xc7, 0xe0, 0x28, 0x94, 0x29, 0x29, 0x49, 0x68, 0x28, 0x39,
	0x02, 0xb3, 0xb3, 0x16, 0x4a, 0x43, 0xf1, 0xe4, 0xf7, 0x66, 0xde, 0x8c, 0xdf, 0x68, 0x01, 0xba,
	0x56, 0xcd, 0xf3, 0x76, 0x65, 0xac, 0x11, 0xc1, 0xb2, 0xac, 0x9b, 0xf4, 0x1e, 0xa2, 0x87, 0x8d,
	0x6a, 0x6c, 0x27, 0xce, 0x20, 0x9c, 0x69, 0x63, 0x96, 0x31, 0x9b, 0xb0, 0xec, 0xbf, 0xf4, 0x42,
	0x5c, 0x42, 0xa4, 0xa8, 0x1f, 0x0f, 0x27, 0x3c, 0x1b, 0xdf, 0x8c, 0x73, 0x37, 0x96, 0xd3, 0x8c,
	0xec, 0x5b, 0xe9, 0x96, 0x41, 0x48, 0x95, 0x3f, 0x96, 0x08, 0x08, 0x16, 0xa5, 0x2d, 0x71, 0x05,
	0xcb, 0x46, 0x92, 0xb8, 0x73, 0xea, 0xba, 0x51, 0x5d, 0xcc, 0xb1, 0x18, 0x4a, 0x2f, 0x9c, 0xb3,
	0x2d, 0xed, 0x63, 0x1c, 0x78, 0xa7, 0xe3, 0x2e, 0x42, 0x55, 0x2b, 0xbd, 0xe8, 0xe2, 0xf0, 0x38,
	0xc2, 0xd4, 0xd5, 0x64, 0xdf, 0x12, 0x27, 0x30, 0xc4, 0x8c, 0x11, 0x8e, 0x71, 0x89, 0x2c, 0xbd,
	0x86, 0x90, 0x0c, 0xe2, 0x14, 0xf8, 0x93, 0x7a, 0xa6, 0x3c, 0x23, 0xe9, 0xa8, 0xfb, 0xf3, 0xa6,
	0xd4, 0x6b, 0xd5, 0xc7, 0xf1, 0x22, 0x2d, 0x60, 0x3c, 0xad, 0xb5, 0x2a, 0x4c, 0x53, 0x5b, 0xb3,
	0xfa, 0x0d, 0xc2, 0x8e, 0x82, 0x9c, 0x43, 0x64, 0xaa, 0xaa, 0x53, 0x96, 0x26, 0xb9, 0xec, 0x15,
	0x79, 0x8d, 0xd6, 0x74, 0xc9, 0x3f, 0x49, 0x3c, 0xbd, 0x80, 0xa0, 0x50, 0xfe, 0xcc, 0xb9, 0x59,
	0x37, 0x96, 0x16, 0x71, 0xe9, 0xc5, 0xdd, 0xd5, 0x6e, 0x9f, 0x0c, 0xde, 0x11, 0xdf, 0xfb, 0x84,
	0xbd, 0x1c, 0x12, 0xf6, 0x8a, 0x78, 0x43, 0xec, 0x10, 0x1f, 0x88, 0xaf, 0x03, 0xf6, 0xf0, 0xbb,
	0xfd, 0x4c, 0x06, 0xb3, 0x88, 0x1e, 0xec, 0xf6, 0x27, 0x00, 0x00, 0xff, 0xff, 0x7f, 0x5a, 0x67,
	0xe8, 0xbe, 0x01, 0x00, 0x00,
}