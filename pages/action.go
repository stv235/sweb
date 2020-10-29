package pages

import "io"

type action interface {}

type redirector struct {
	url string
}

type renderer struct {
	name string
	data interface{}
}

type file struct {
	buf      []byte
	mimeType string
}

type status struct {
	code int
	message string
}

type send struct {
	mimeType string
	r        io.Reader
	closeFn  func()
}

type requireAuth struct {
	realm string
}