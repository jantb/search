package main

import (
	"io"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

type Field struct {
	Key   string
	Value string
}

func (d *Field) Size() (s uint64) {

	{
		l := uint64(len(d.Key))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Value))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *Field) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Key))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Key)
		i += l
	}
	{
		l := uint64(len(d.Value))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Value)
		i += l
	}
	return buf[:i+0], nil
}

func (d *Field) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Key = string(buf[i+0 : i+0+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Value = string(buf[i+0 : i+0+l])
		i += l
	}
	return i + 0, nil
}

type Event struct {
	Bloom  []byte
	Ts     time.Time
	Data   string
	Lines  int32
	Path   string
	Fields []Field
}

func (d *Event) Size() (s uint64) {

	{
		l := uint64(len(d.Bloom))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Data))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Path))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Fields))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}

		for k := range d.Fields {

			{
				s += d.Fields[k].Size()
			}

		}

	}
	s += 19
	return
}
func (d *Event) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Bloom))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i:], d.Bloom)
		i += l
	}
	{
		b, err := d.Ts.MarshalBinary()
		if err != nil {
			return nil, err
		}
		copy(buf[i+0:], b)
	}
	{
		l := uint64(len(d.Data))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+15] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+15] = byte(t)
			i++

		}
		copy(buf[i+15:], d.Data)
		i += l
	}
	{

		buf[i+0+15] = byte(d.Lines >> 0)

		buf[i+1+15] = byte(d.Lines >> 8)

		buf[i+2+15] = byte(d.Lines >> 16)

		buf[i+3+15] = byte(d.Lines >> 24)

	}
	{
		l := uint64(len(d.Path))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+19] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+19] = byte(t)
			i++

		}
		copy(buf[i+19:], d.Path)
		i += l
	}
	{
		l := uint64(len(d.Fields))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+19] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+19] = byte(t)
			i++

		}
		for k := range d.Fields {

			{
				nbuf, err := d.Fields[k].Marshal(buf[i+19:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		}
	}
	return buf[:i+19], nil
}

func (d *Event) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Bloom)) >= l {
			d.Bloom = d.Bloom[:l]
		} else {
			d.Bloom = make([]byte, l)
		}
		copy(d.Bloom, buf[i:])
		i += l
	}
	{
		d.Ts.UnmarshalBinary(buf[i+0 : i+0+15])
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+15] & 0x7F)
			for buf[i+15]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+15]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Data = string(buf[i+15 : i+15+l])
		i += l
	}
	{

		d.Lines = 0 | (int32(buf[i+0+15]) << 0) | (int32(buf[i+1+15]) << 8) | (int32(buf[i+2+15]) << 16) | (int32(buf[i+3+15]) << 24)

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+19] & 0x7F)
			for buf[i+19]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+19]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Path = string(buf[i+19 : i+19+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+19] & 0x7F)
			for buf[i+19]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+19]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Fields)) >= l {
			d.Fields = d.Fields[:l]
		} else {
			d.Fields = make([]Field, l)
		}
		for k := range d.Fields {

			{
				ni, err := d.Fields[k].Unmarshal(buf[i+19:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

		}
	}
	return i + 19, nil
}
