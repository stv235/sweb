package randutil

import (
	"crypto/rand"
	"log"
)

func NewRandomBuffer(length int) []byte {
	buf := make([]byte, length)

	_, err := rand.Read(buf)

	if err != nil {
		log.Panicln(err)
	}

	return buf
}
