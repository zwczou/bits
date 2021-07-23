package bits

import "io"

type bufReader struct {
	buf []byte
	idx int
}

func (r *bufReader) Read(b []byte) (n int, err error) {
	// 长度不够报错, 证明解包有问题
	if len(b) > len(r.buf)-r.idx {
		err = io.ErrUnexpectedEOF
		return
	}

	copy(b, r.buf[r.idx:r.idx+len(b)])
	r.idx += len(b)
	return len(b), nil
}

type bufWriter struct {
	buf []byte
	idx int
}

func (w *bufWriter) Write(b []byte) (n int, err error) {
	// 如果长度不够，不是写入剩下长度，而是直接报错
	// 证明上级分配长度不够
	if len(b) > len(w.buf)-w.idx {
		err = io.ErrUnexpectedEOF
		return
	}

	copy(w.buf[w.idx:], b)
	w.idx += len(b)
	return len(b), nil
}
