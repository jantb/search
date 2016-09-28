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

type FileMonitor struct {
	Path   string
	Offset int64
	Poll   bool
}

func (d *FileMonitor) Size() (s uint64) {

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
	s += 9
	return
}
func (d *FileMonitor) Marshal(buf []byte) ([]byte, error) {
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
		l := uint64(len(d.Path))

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
		copy(buf[i+0:], d.Path)
		i += l
	}
	{

		buf[i+0+0] = byte(d.Offset >> 0)

		buf[i+1+0] = byte(d.Offset >> 8)

		buf[i+2+0] = byte(d.Offset >> 16)

		buf[i+3+0] = byte(d.Offset >> 24)

		buf[i+4+0] = byte(d.Offset >> 32)

		buf[i+5+0] = byte(d.Offset >> 40)

		buf[i+6+0] = byte(d.Offset >> 48)

		buf[i+7+0] = byte(d.Offset >> 56)

	}
	{
		if d.Poll {
			buf[i+8] = 1
		} else {
			buf[i+8] = 0
		}
	}
	return buf[:i+9], nil
}

func (d *FileMonitor) Unmarshal(buf []byte) (uint64, error) {
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
		d.Path = string(buf[i+0 : i+0+l])
		i += l
	}
	{

		d.Offset = 0 | (int64(buf[i+0+0]) << 0) | (int64(buf[i+1+0]) << 8) | (int64(buf[i+2+0]) << 16) | (int64(buf[i+3+0]) << 24) | (int64(buf[i+4+0]) << 32) | (int64(buf[i+5+0]) << 40) | (int64(buf[i+6+0]) << 48) | (int64(buf[i+7+0]) << 56)

	}
	{
		d.Poll = buf[i+8] == 1
	}
	return i + 9, nil
}
