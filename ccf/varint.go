package ccf

import (
	"encoding/binary"
	"io"
)

func writeUvarint(w io.Writer, v uint64) error {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, v)

	if _, err := w.Write(buf[:n]); err != nil {
		return err
	}

	return nil
}

func writeVarint(w io.Writer, v int64) error {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, v)

	if _, err := w.Write(buf[:n]); err != nil {
		return err
	}

	return nil
}
