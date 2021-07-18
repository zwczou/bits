package bits

import (
	"io"
)

type Reader struct {
	br         *BitReader
	err        error
	checkError CheckError
}

func NewReader(ir io.Reader) *Reader {
	return &Reader{
		br:         NewBitReader(ir),
		checkError: Check(),
	}
}

func NewReaderBuffer(buf []byte) *Reader {
	return &Reader{
		br:         NewBitReaderBuffer(buf),
		checkError: Check(),
	}
}

// 错误直接panic
func (r *Reader) Must() *Reader {
	r.checkError = Must()
	return r
}

func (r *Reader) Error() error {
	return r.err
}

func (r *Reader) Index() int {
	return r.br.Index()
}

func (r *Reader) Reset(ir io.Reader) {
	r.br.Reset(ir)
}

func (r *Reader) Read(b []byte) int {
	n, err := r.br.Read(b)
	r.err = r.checkError(err)
	return n
}

func (r *Reader) ReadBool() bool {
	b, err := r.br.ReadBit()
	r.err = r.checkError(err)
	return bool(b)
}

func (r *Reader) ReadUint8(nbits int) uint8 {
	if nbits > 8 {
		r.err = r.checkError(ErrOverflow)
		return 0
	}

	u, err := r.br.ReadBits(nbits)
	r.err = r.checkError(err)
	return uint8(u)
}

func (r *Reader) ReadUint16(nbits int) uint16 {
	if nbits > 16 {
		r.err = r.checkError(ErrOverflow)
		return 0
	}

	u, err := r.br.ReadBits(nbits)
	r.err = r.checkError(err)
	return uint16(u)
}

func (r *Reader) ReadUint32(nbits int) uint32 {
	if nbits > 32 {
		r.err = r.checkError(ErrOverflow)
		return 0
	}

	u, err := r.br.ReadBits(nbits)
	r.err = r.checkError(err)
	return uint32(u)
}

func (r *Reader) ReadUint64(nbits int) uint64 {
	if nbits > 64 {
		r.err = r.checkError(ErrOverflow)
		return 0
	}

	u, err := r.br.ReadBits(nbits)
	r.err = r.checkError(err)
	return u
}

func (r *Reader) Skip(nbits int) {
	r.err = r.checkError(r.br.Skip(nbits))
}

func (r *Reader) Align() {
	r.br.Align()
}

func (r *Reader) IsAligned() bool {
	return r.br.IsAligned()
}