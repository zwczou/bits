package bits

import "io"

type bufReader struct {
	buf []byte
	idx int
}

func (r *bufReader) Read(b []byte) (n int, err error) {
	size := len(r.buf) - r.idx
	if size <= 0 {
		return 0, io.EOF
	}

	if len(b) < size {
		size = len(b)
	}

	copy(b, r.buf[r.idx:r.idx+size])
	r.idx += size
	return size, nil
}

type bufWriter struct {
	buf []byte
	idx int
}

func (w *bufWriter) Write(b []byte) (n int, err error) {
	// 如果长度不够，不是写入剩下长度，而是直接报错
	// 证明上级分配长度不够
	size := len(w.buf) - w.idx
	if size <= 0 {
		return 0, io.EOF
	}
	if len(b) < size {
		size = len(b)
	}

	copy(w.buf[w.idx:], b[:size])
	w.idx += size
	return size, nil
}
