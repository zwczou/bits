package bits

import (
	"bytes"
	"io"
	"testing"
)

func TestBitReader(t *testing.T) {
	buf := []byte{1, 2, 3, 7}
	r := NewBitReaderBuffer(buf)
	for i := 0; i < 7; i++ {
		b, err := r.ReadBit()
		if err != nil {
			t.Fatal(err)
		}
		if b {
			t.FailNow()
		}
	}
	if r.Index() != 1 {
		t.FailNow()
	}
	b, err := r.ReadBit()
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.FailNow()
	}

	val, err := r.ReadBits(7)
	if err != nil {
		t.Fatal(err)
	}
	if val != 1 {
		t.FailNow()
	}
	val, err = r.ReadBits(1)
	if err != nil {
		t.Fatal(err)
	}
	if val != 0 {
		t.FailNow()
	}
	val, err = r.ReadBits(16)
	if err != nil {
		t.Fatal(err)
	}
	if val != 0x307 {
		t.FailNow()
	}
	if r.Index() != 4 {
		t.FailNow()
	}
	_, err = r.ReadBit()
	if err != io.ErrUnexpectedEOF {
		t.Fatal(err)
	}

	r.Reset(bytes.NewReader(buf))
	_, err = r.ReadBit()
	if err != nil {
		t.Fatal(err)
	}
	err = r.Skip(27)
	if err != nil {
		t.Fatal(err)
	}
	b, err = r.ReadBit()
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.FailNow()
	}
	u, err := r.ReadBits(3)
	if err != nil {
		t.Fatal(err)
	}
	if u != 0x7 {
		t.FailNow()
	}
}

func TestBitWriter(t *testing.T) {
	buf := make([]byte, 8)
	w := NewBitWriterBuffer(buf)
	for i := 0; i < 8; i++ {
		if (i+1)%4 == 0 {
			w.WriteBit(true)
		} else {
			w.WriteBit(false)
		}
	}

	w.WriteByte(0x1f)
	w.WriteBit(false)
	w.WriteBit(false)
	w.WriteBit(false)
	w.WriteBit(true)
	w.WriteBits(0xf, 4)

	w.WriteBits(0x1f, 3)
	w.Flush(true) // padding 0b11111
	if w.Index() != 4 {
		t.FailNow()
	}

	w.WriteBits(0x11223344, 32)

	if buf[0] != 0x11 || buf[1] != 0x1f || buf[2] != 0x1f || buf[3] != 0xff {
		t.FailNow()
	}
	if buf[4] != 0x11 || buf[5] != 0x22 || buf[6] != 0x33 || buf[7] != 0x44 {
		t.FailNow()
	}
	if w.Index() != 8 {
		t.FailNow()
	}
	err := w.WriteByte(1)
	if err != io.ErrUnexpectedEOF {
		t.FailNow()
	}
}
