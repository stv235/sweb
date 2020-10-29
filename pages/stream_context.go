package pages

import (
	"io"
	"log"
	"mime"
	"mime/multipart"
)

type StreamContext struct {
	CommonContext

	ContentType string

	W    io.Writer
	R    io.ReadCloser
}

func (ctx *StreamContext) ReadMultipartForm() *multipart.Form {
	_, params, err := mime.ParseMediaType(ctx.ContentType)

	if err != nil {
		log.Panicln(err)
	}

	b, ok := params["boundary"]

	if !ok {
		log.Panicln("[HTTP]", "no boundary found")
	}

	r := multipart.NewReader(ctx.R, b)

	f, err := r.ReadForm(10 * 1024 * 1024)

	if err != nil {
		log.Panicln(err)
	}

	return f
}

func (ctx *StreamContext) Send(r io.Reader, mimeType string, closeFn func()) {
	ctx.action = &send{ r: r, mimeType: mimeType, closeFn: closeFn }
}

