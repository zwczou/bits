// 这部分代码主要来源于此
// https://github.com/dgryski/go-bitstream/blob/master/bitstream.go
package bits

import (
	"bytes"
	"io"
)

type Bit bool

const (
	Zero Bit = false
	One  Bit = true
)

type BitReader struct {
	r     io.Reader
	b     [1]byte
	count uint8
	idx   int
}

func NewBitReader(r io.Reader) *BitReader {
	return &BitReader{
		r: r,
	}
}

func NewBitReaderBuffer(buf []byte) *BitReader {
	return &BitReader{
		r: bytes.NewReader(buf),
	}
}

func (br *BitReader) Index() int {
	return br.idx
}

func (br *BitReader) SetIndex(idx int) {
	br.idx = idx
}

func (br *BitReader) Reset(ir io.Reader) {
	br.r = ir
	br.b[0] = 0x00
	br.idx = 0
	br.count = 0
}

func (br *BitReader) ReadBit() (bit Bit, err error) {
	if br.count == 0 {
		_, err = br.r.Read(br.b[:])
		if err != nil {
			return
		}
		br.idx++
		br.count = 8
	}

	br.count--
	d := (br.b[0] & 0x80)
	br.b[0] <<= 1
	return d != 0, nil
}

func (br *BitReader) ReadByte() (byt byte, err error) {
	if br.count == 0 {
		_, err = br.r.Read(br.b[:])
		if err != nil {
			return
		}
		br.idx++
		return br.b[0], nil
	}

	byt = br.b[0]
	_, err = br.r.Read(br.b[:])
	if err != nil {
		return
	}
	br.idx++

	byt |= br.b[0] >> br.count
	br.b[0] <<= (8 - br.count)
	return byt, nil
}

func (br *BitReader) ReadBits(nbits int) (u uint64, err error) {
	for nbits >= 8 {
		byt, err := br.ReadByte()
		if err != nil {
			return 0, err
		}

		u = (u << 8) | uint64(byt)
		nbits -= 8
	}

	for nbits > 0 {
		byt, err := br.ReadBit()
		if err != nil {
			return 0, err
		}
		u <<= 1
		if byt {
			u |= 1
		}
		nbits--
	}
	return
}

// 性能并不是很高
func (br *BitReader) Read(b []byte) (n int, err error) {
	if br.count == 0 {
		n, err = io.ReadFull(br.r, b)
		br.idx += n
		return
	}

	for i := 0; i < len(b); i++ {
		b[i], err = br.ReadByte()
		if err != nil {
			return
		}
		n += 1
	}
	return
}

func (br *BitReader) Skip(nbits int) (err error) {
	// 首先对齐
	for nbits > 0 && br.count > 0 {
		br.count--
		br.b[0] <<= 1
		nbits--
	}

	// 往上多取一个字节
	want := (nbits + 7) / 8
	if want > 0 {
		b := make([]byte, want)
		_, err = io.ReadFull(br.r, b)
		if err != nil {
			return
		}
		br.idx += want

		br.count = 8
		br.b[0] = b[want-1]
		nbits -= (want - 1) * 8
	}

	if nbits > 0 {
		br.count -= uint8(nbits)
		br.b[0] <<= nbits
	}
	return
}

func (br *BitReader) IsAligned() bool {
	return br.count == 0
}

func (br *BitReader) Align() {
	br.Skip(int(br.count))
}

type bufWriter struct {
	buf []byte
	idx int
}

func (w *bufWriter) Write(b []byte) (n int, err error) {
	if len(b) > len(w.buf)-w.idx {
		err = io.ErrUnexpectedEOF
		return
	}
	for _, c := range b {
		w.buf[w.idx] = c
		w.idx++
	}
	return len(b), nil
}

type BitWriter struct {
	w     io.Writer
	b     [1]byte
	count uint8
	idx   int
}

func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{
		w:     w,
		count: 8,
	}
}

func NewBitWriterBuffer(buf []byte) *BitWriter {
	return &BitWriter{
		w: &bufWriter{
			buf: buf,
		},
		count: 8,
	}
}

func (bw *BitWriter) Index() int {
	return bw.idx
}

func (bw *BitWriter) SetIndex(idx int) {
	bw.idx = idx
}

func (bw *BitWriter) Resume(data byte, count uint8) {
	bw.b[0] = data
	bw.count = count
}

func (bw *BitWriter) Reset(iw io.Writer) {
	bw.w = iw
	bw.b[0] = 0x00
	bw.idx = 0
	bw.count = 8
}

func (bw *BitWriter) WriteBit(bit Bit) (err error) {
	if bit {
		bw.b[0] |= 1 << (bw.count - 1)
	}

	bw.count--

	if bw.count == 0 {
		_, err = bw.w.Write(bw.b[:])
		if err != nil {
			return
		}
		bw.idx++
		bw.b[0] = 0
		bw.count = 8
	}
	return
}

func (bw *BitWriter) WriteByte(byt byte) (err error) {
	bw.b[0] |= byt >> (8 - bw.count)

	_, err = bw.w.Write(bw.b[:])
	if err != nil {
		return
	}
	bw.idx++

	bw.b[0] = byt << bw.count
	return
}

func (bw *BitWriter) WriteBits(u uint64, nbits int) (err error) {
	u <<= (64 - uint(nbits))

	for nbits >= 8 {
		byt := byte(u >> 56)
		err = bw.WriteByte(byt)
		if err != nil {
			return
		}
		u <<= 8
		nbits -= 8
	}

	for nbits > 0 {
		err = bw.WriteBit((u >> 63) == 1)
		if err != nil {
			return
		}
		u <<= 1
		nbits--
	}
	return
}

func (bw *BitWriter) Write(b []byte) (n int, err error) {
	if bw.count == 8 {
		n, err = bw.w.Write(b)
		bw.idx += n
		return
	}

	for _, c := range b {
		err = bw.WriteByte(c)
		if err != nil {
			return
		}
		n += 1
	}
	return
}

func (bw *BitWriter) Flush(bit Bit) (err error) {
	for bw.count != 8 {
		err = bw.WriteBit(bit)
		if err != nil {
			return
		}
	}
	return
}

func (bw *BitWriter) IsAligned() bool {
	return bw.count == 8
}
