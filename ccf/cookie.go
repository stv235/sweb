package ccf

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"reflect"
	"time"
)

// compact cookie format

var ErrInvalidType = errors.New("invalid type")
var ErrInvalidTag = errors.New("invalid tag")
var ErrTooShort = errors.New("buffer too short")
var ErrInvalidHash = errors.New("invalid hash")
var ErrInvalidFormat = errors.New("invalid format")
var ErrTimeout = errors.New("timeout")


type Cookie struct {
	buf []byte

}

func createId(name string) (Id, error) {
	if len(name) != 1 {
		return 0, ErrInvalidTag
	}

	return Id(name[0]), nil
}

func (cookie Cookie) Unmarshal(s interface{}) error {
	r := bytes.NewReader(cookie.buf)

	fields := make(map[Id]reflect.Value)

	v := reflect.ValueOf(s)
	v = reflect.Indirect(v)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)

		if name, ok := ft.Tag.Lookup("ccf"); ok {
			id, err := createId(name)

			if err != nil {
				return err
			}

			fields[id] = v.Field(i)
		}
	}

	for r.Len() > 0 {
		id, err := r.ReadByte()

		if err != nil {
			return err
		}

		if f, ok := fields[Id(id)]; ok {
			switch f.Interface().(type) {
			case string:
				buf, err := decodeBytes(r)

				if err != nil {
					return err
				}

				f.SetString(string(buf))
			case time.Time:
				x, err := binary.ReadVarint(r)

				if err != nil {
					return err
				}
				f.Set(reflect.ValueOf(time.Unix(x / 1000000000, x % 1000000000)))
			case int64:
				x, err := binary.ReadVarint(r)

				if err != nil {
					return err
				}

				f.SetInt(x)
			case *int64:
				x, err := binary.ReadVarint(r)

				if err != nil {
					return err
				}

				f.Set(reflect.ValueOf(&x))
			}
		}
	}

	return nil
}

func (cookie *Cookie) Marshal(s interface{}) error {
	b := bytes.NewBuffer(nil)

	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)

		if name, ok := ft.Tag.Lookup("ccf"); ok {
			id, err := createId(name)

			if err != nil {
				return err
			}

			f := v.Field(i)

			if err := encodeValue(b, id, f.Interface()); err != nil {
				return err
			}
		}
	}

	cookie.buf = b.Bytes()

	return nil
}

func encodeValue(b *bytes.Buffer, id Id, val interface{}) error {
	// allow optional values
	switch val.(type) {
	case *int64:
		if val.(*int64) != nil {
			return encodeValue(b, id, *(val.(*int64)))
		}

		return nil
	}

	b.WriteByte(byte(id))

	switch val.(type) {
	case string:
		return encodeBytes(b, []byte(val.(string)))
	case int64:
		return writeVarint(b, val.(int64))
	case time.Time:
		return writeVarint(b, val.(time.Time).UTC().UnixNano())
	}

	return ErrInvalidType
}

func encodeBytes(b *bytes.Buffer, buf []byte) error {
	if err := writeUvarint(b, uint64(len(buf))); err != nil {
		return err
	}

	if _, err :=b.Write(buf); err != nil {
		return err
	}

	return nil
}

func decodeBytes(r *bytes.Reader) ([]byte, error) {
	n, err := binary.ReadUvarint(r)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, n)

	if n == 0 {
		return buf, nil
	}

	if _, err := r.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func (cookie Cookie) Encode(timeout time.Time, encryptKey, hashKey string) (string, error) {
	b := bytes.NewBuffer(nil)

	if err := writeVarint(b, timeout.UnixNano()); err != nil {
		return "", err
	}

	b.Write(cookie.buf)

	buf := b.Bytes()

	mac := hmac.New(sha512.New, []byte(hashKey))
	mac.Write(buf)
	hash := mac.Sum(nil)

	buf = append(hash, buf...)

	buf, err := deflate(buf)

	if err != nil {
		return "", err
	}

	buf, err = encrypt(buf, encryptKey)

	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(buf), nil
}

func (cookie *Cookie) Decode(str, encryptKey, hashKey string) error {
	buf, err := base64.RawStdEncoding.DecodeString(str)

	if err != nil {
		return err
	}

	buf, err = decrypt(buf, encryptKey)

	if err != nil {
		return err
	}

	buf, err = inflate(buf)

	if err != nil {
		return err
	}

	if len(buf) < sha512.Size {
		return ErrTooShort
	}

	mac := hmac.New(sha512.New, []byte(hashKey))
	requiredHash := buf[:sha512.Size]
	buf = buf[sha512.Size:]

	mac.Write(buf)
	calculatedHash := mac.Sum(nil)

	if !hmac.Equal(calculatedHash, requiredHash) {
		return ErrInvalidHash
	}

	timestamp, n := binary.Varint(buf)

	if n <= 0 {
		return ErrInvalidFormat
	}

	timeout := time.Unix(timestamp / 1000000000, timestamp % 1000000000)

	if timeout.Before(time.Now()) {
		return ErrTimeout
	}

	cookie.buf = buf[n:]

	return nil
}