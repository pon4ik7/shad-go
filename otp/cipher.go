//go:build !solution

package otp

import (
	"io"
)

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &xorReader{reader: r, prng: prng}
}

type xorReader struct {
	reader io.Reader
	prng   io.Reader
}

func (x *xorReader) Read(p []byte) (int, error) {
	n, err := x.reader.Read(p)
	ks := make([]byte, n)
	if n > 0 {
		m, _ := x.prng.Read(ks)
		n = min(n, m)
		for i := 0; i < n; i++ {
			p[i] ^= ks[i]
		}
	}
	return n, err
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &xorWriter{writer: w, prng: prng}
}

type xorWriter struct {
	writer io.Writer
	prng   io.Reader
}

func (x *xorWriter) Write(p []byte) (int, error) {
	buf := make([]byte, len(p))
	m, _ := x.prng.Read(buf)
	n := min(len(p), m)
	for i := 0; i < n; i++ {
		buf[i] ^= p[i]
	}
	n, err := x.writer.Write(buf[:n])
	return n, err
}
