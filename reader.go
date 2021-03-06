package bits

import (
	"io"
)

type Reader struct {
	br         *BitReader
	err        error
	checkError CheckError
}

// 默认出错直接panic，所以需要在上层做一些长度等判断
// 如果需要自己判断错误，则用
// br := NewReader(iw).Check()
// br.ReadBool()
// br.Error() // 检查ReadBool()错误
func NewReader(ir io.Reader) *Reader {
	return &Reader{
		br:         NewBitReader(ir),
		checkError: Must(),
	}
}

func NewReaderBuffer(buf []byte) *Reader {
	return &Reader{
		br:         NewBitReader(NewBufferReader(buf)),
		checkError: Must(),
	}
}

// 错误直接panic
func (r *Reader) Must() *Reader {
	r.checkError = Must()
	return r
}

// 需要每次判断r.Error()
func (r *Reader) Check() *Reader {
	r.checkError = Check()
	return r
}

func (r *Reader) Error() error {
	return r.err
}

// 只能保留最后一个错误
func (r *Reader) SetError(err error) {
	err = r.checkError(err)
	if err != nil {
		r.err = err
	}
}

func (r *Reader) Index() int {
	return r.br.Index()
}

func (r *Reader) Reset(ir io.Reader) {
	r.br.Reset(ir)
}

func (r *Reader) Reader() io.Reader {
	return r.br.Reader()
}

func (r *Reader) Read(b []byte) int {
	n, err := r.br.Read(b)
	r.SetError(err)
	return n
}

func (r *Reader) ReadBool() bool {
	b, err := r.br.ReadBit()
	r.SetError(err)
	return bool(b)
}

func (r *Reader) ReadUint8(nbits int) uint8 {
	if nbits > 8 {
		r.SetError(ErrOverflow)
		return 0
	}

	u, err := r.br.ReadBits(nbits)
	r.SetError(err)
	return uint8(u)
}

func (r *Reader) ReadUint16(nbits int) uint16 {
	if nbits > 16 {
		r.SetError(ErrOverflow)
		return 0
	}

	u, err := r.br.ReadBits(nbits)
	r.SetError(err)
	return uint16(u)
}

func (r *Reader) ReadUint32(nbits int) uint32 {
	if nbits > 32 {
		r.SetError(ErrOverflow)
		return 0
	}

	u, err := r.br.ReadBits(nbits)
	r.SetError(err)
	return uint32(u)
}

func (r *Reader) ReadUint64(nbits int) uint64 {
	if nbits > 64 {
		r.SetError(ErrOverflow)
		return 0
	}

	u, err := r.br.ReadBits(nbits)
	r.SetError(err)
	return u
}

func (r *Reader) Skip(nbits int) {
	r.SetError(r.br.Skip(nbits))
}

func (r *Reader) Align() {
	r.br.Align()
}

func (r *Reader) IsAligned() bool {
	return r.br.IsAligned()
}
