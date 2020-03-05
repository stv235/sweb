package randutil

import (
	"encoding/base64"
)

func NewRandomString(length int) string {
	return base64.RawStdEncoding.EncodeToString(NewRandomBuffer(length))
}