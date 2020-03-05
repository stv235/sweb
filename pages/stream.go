package pages

import (
	"sweb/ccf"
	"net/http"
	"net/url"
)

type Stream struct {
	Get func(ctx *StreamContext)
	Post func(ctx *StreamContext)
}

func (stream Stream) handle(w http.ResponseWriter, r *http.Request, pathParameters url.Values, sessionParameters *ccf.Cookie) (action, error) {
	context := StreamContext{}
	context.Path = r.URL.Path
	context.PathParameters = pathParameters
	context.ContentType = r.Header.Get("Content-Type")

	context.W = w

	switch {
	case r.Method == http.MethodGet && stream.Get != nil:
		stream.Get(&context)
	case r.Method == http.MethodPost && stream.Post != nil:
		context.R = r.Body
		stream.Post(&context)
	default:
		return nil, ErrMethodNotAllowed
	}

	return context.action, nil
}

