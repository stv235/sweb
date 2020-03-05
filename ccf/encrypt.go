package ccf

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
)

func createKey(password string) []byte {
	h := sha256.New()
	h.Write([]byte(password))

	return h.Sum(nil)
}

func makeIv(block cipher.Block) ([]byte, error) {
	iv := make([]byte, block.BlockSize())

	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	return iv, nil
}

func encrypt(buf []byte, password string) ([]byte, error) {
	block, err := aes.NewCipher(createKey(password))

	if err != nil {
		return nil, err
	}

	iv, err := makeIv(block)

	if err != nil {
		return nil, err
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(buf, buf)

	return append(iv, buf...), nil
}

func decrypt(buf []byte, password string) ([]byte, error) {
	block, err := aes.NewCipher(createKey(password))

	if err != nil {
		return nil, err
	}

	iv := buf[:block.BlockSize()]
	buf = buf[block.BlockSize():]

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(buf, buf)

	return buf, nil
}
