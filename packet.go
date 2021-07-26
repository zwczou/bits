package bits

import "io"

type Packet interface {
	io.Reader
	io.Writer
	Size() int
	ReadBits(*Reader)
	WriteBits(*Writer)
}
