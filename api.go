package bits

import "io"

type Packer interface {
	io.Reader
	ReadBits(*Reader)
}

type Unpacker interface {
	io.Writer
	WriteBits(*Writer)
}

type Sizer interface {
	Size() int
}

type Packeter interface {
	Sizer
	Packer
	Unpacker
}
