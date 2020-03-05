package ccf

import (
	"bytes"
	"compress/flate"
	"io"
	"io/ioutil"
)

func inflate(buf []byte) ([]byte, error) {
	r1 := bytes.NewReader(buf)
	r2 := flate.NewReader(r1)

	buf, err := ioutil.ReadAll(r2)

	if err != nil && err != io.EOF {
		return nil, err
	}

	return buf, nil
}

func deflate(buf []byte) ([]byte, error) {
	w1 := bytes.NewBuffer(nil)
	w2, err := flate.NewWriter(w1, flate.BestCompression)

	if err != nil {
		return nil, err
	}

	if _, err := w2.Write(buf); err != nil {
		return nil, err
	}

	if err := w2.Close(); err != nil {
		return nil, err
	}

	return w1.Bytes(), nil
}
