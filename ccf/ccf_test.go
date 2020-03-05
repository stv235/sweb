package ccf

import (
	"log"
	"strings"
	"testing"
	"time"
)

type TestData struct {
	Str string `ccf:"s"`
	Int int64 `ccf:"i"`
	T time.Time `ccf:"t"`
}

func TestEncode(t *testing.T) {
	d := TestData{}
	d.T = time.Now()
	d.Int = 100
	d.Str = "200"

	buf, err := Encode(d)

	if err != nil {
		t.Fail()
	}

	log.Println(string(buf))

	d2 := TestData{}

	if err := Decode(buf, &d2); err != nil {
		log.Println(err)
		t.Fail()
	}

	if d.Str != d2.Str {
		log.Println("str != str")
		t.Fail()
	}

	if d.Int != d2.Int {
		log.Println("int != int")
		t.Fail()
	}

	if d.T.UTC().UnixNano() != d2.T.UTC().UnixNano() {
		log.Println("T != T")
		log.Println(d.T)
		log.Println(d2.T)
		t.Fail()
	}
}

func TestCookie_Encode(t *testing.T) {
	d1 := TestData{}
	d1.Str = strings.Repeat("abc", 1024)
	d1.Int = 9999
	d1.T = time.Now()

	hashKey := "xyz"
	encryptKey := "abc"

	cookie := Cookie{}



	if err := cookie.Marshal(d1); err != nil {
		log.Println(err)
		t.Fail()
	}

	str, err := cookie.Encode(time.Now().Add(time.Minute * 5), encryptKey, hashKey)

	if err != nil {
		t.Fail()
	}

	err = cookie.Decode(str, encryptKey, hashKey)

	if err != nil {
		t.Fail()
	}

	d2 := TestData{}
	if err := cookie.Unmarshal(&d2); err != nil {
		log.Println(err)
		t.Fail()
	}

	if d1.Str != d2.Str {
		log.Println("d1.Str != d2.Str")
		t.Fail()
	}
}