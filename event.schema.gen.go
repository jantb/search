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

type Event struct {
	Bloom []byte
	Ts    time.Time
	Data  string
	Lines int32
	Path  string
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
	return i + 19, nil
}
