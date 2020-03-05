package password

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"log"
)

func createHash(buf []byte, key []byte) []byte {
	if len(key) == 0 {
		log.Panicln("empty hash key")
	}

	mac := hmac.New(sha512.New, key)
	mac.Write(buf)

	return mac.Sum(nil)
}

func Create(buf []byte, key string) string {
	hash := createHash(buf, []byte(key))

	return base64.StdEncoding.EncodeToString(hash)
}

func Compare(buf []byte, hashStr string, key string) (bool, error) {
	hash1, err := base64.StdEncoding.DecodeString(hashStr)

	if err != nil {
		return false, err
	}

	hash2 := createHash(buf, []byte(key))

	return hmac.Equal(hash1, hash2), nil
}