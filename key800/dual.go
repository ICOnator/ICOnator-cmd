package main

import (
	"io"
	"errors"
)

type dualReader struct {
	reader1 io.Reader
	reader2 io.Reader
}

func New(reader1 io.Reader, reader2 io.Reader) *dualReader {
	return &dualReader{reader1, reader2}
}

func (mr *dualReader) Read(p []byte) (n int, err error) {

	tmp1 := make([]byte, len(p))
	tmp2 := make([]byte, len(p))

	n1, err1 := mr.reader1.Read(tmp1)
	n2, err2 := mr.reader2.Read(tmp2)

	if err1 != nil {
		return n1, err1;
	}

	if err2 != nil {
		return n2, err2;
	}

	if len(p) != n1 && len(p) != n2 {
		return 0, errors.New("did not read same length")
	}

	for i := 0; i<n1; i++ {
		p[i] = tmp1[i] ^ tmp2[i];
	}

	return n1, nil
}