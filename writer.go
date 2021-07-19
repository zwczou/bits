package bits

import "io"

type Writer struct {
	bw         *BitWriter
	checkError CheckError
	err        error
}

// 默认出错直接panic，所以需要在上层做一些长度等判断
// 如果需要自己判断错误，则用
// bw := NewWriter(iw).Check() // 初始化
// bw.WriteBool(false)
// bw.Error() // 检查bw.WriteBool出错
func NewWriter(iw io.Writer) *Writer {
	return &Writer{
		bw:         NewBitWriter(iw),
		checkError: Check(),
	}
}

func NewWriterBuffer(buf []byte) *Writer {
	return &Writer{
		bw:         NewBitWriterBuffer(buf),
		checkError: Must(),
	}
}

func (w *Writer) Must() *Writer {
	w.checkError = Must()
	return w
}

// 需要每次判断w.Error()
func (w *Writer) Check() *Writer {
	w.checkError = Check()
	return w
}

func (w *Writer) Error() error {
	return w.err
}

func (w *Writer) Index() int {
	return w.bw.Index()
}

func (w *Writer) Reset(iw io.Writer) {
	w.bw.Reset(iw)
}

func (w *Writer) Resume(data byte, count uint8) {
	w.bw.Resume(data, count)
}

func (w *Writer) Write(b []byte) (n int) {
	n, err := w.bw.Write(b)
	w.err = w.checkError(err)
	return n
}

func (w *Writer) WriteBool(v bool) {
	w.err = w.checkError(w.bw.WriteBit(Bit(v)))
}

func (w *Writer) WriteUint8(u uint8, nbits int) {
	if nbits > 8 {
		w.err = w.checkError(ErrOverflow)
		return
	}

	w.err = w.checkError(w.bw.WriteBits(uint64(u), nbits))
}

func (w *Writer) WriteUint16(u uint16, nbits int) {
	if nbits > 16 {
		w.err = w.checkError(ErrOverflow)
		return
	}

	w.err = w.checkError(w.bw.WriteBits(uint64(u), nbits))
}

func (w *Writer) WriteUint32(u uint32, nbits int) {
	if nbits > 32 {
		w.err = w.checkError(ErrOverflow)
		return
	}

	w.err = w.checkError(w.bw.WriteBits(uint64(u), nbits))
}

func (w *Writer) WriteUint64(u uint64, nbits int) {
	if nbits > 64 {
		w.err = w.checkError(ErrOverflow)
		return
	}

	w.err = w.checkError(w.bw.WriteBits(u, nbits))
}

func (w *Writer) Flush(bit Bit) {
	w.err = w.checkError(w.bw.Flush(bit))
}

func (w *Writer) IsAligned() bool {
	return w.bw.IsAligned()
}
