package pages

import (
	"net/http"
	"net/url"
	"sweb/ccf"
)

type Page struct {
	Get       func(ctx *Context)
	Post      func(ctx *Context)
}

func (page Page) handleGet(w http.ResponseWriter, r *http.Request, ctx *Context) (action, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	ctx.GetParameters = r.Form
	page.Get(ctx)

	return ctx.action, nil
}

func (page Page) handlePost(w http.ResponseWriter, r *http.Request, ctx *Context) (action, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	ctx.PostParameters = r.Form
	page.Post(ctx)

	return ctx.action, nil
}

func (page Page) handle(w http.ResponseWriter, r *http.Request, pathParameters url.Values, sessionParameters *ccf.Cookie) (action, error) {
	context := Context{}
	context.Path = r.URL.Path
	context.PathParameters = pathParameters
	context.SessionParameters = sessionParameters

	username, password, ok := r.BasicAuth()

	if ok {
		context.Auth = &BasicAuth{
			Username: username,
			Password: password,
		}
	}

	switch {
	case r.Method == http.MethodGet && page.Get != nil:
		return page.handleGet(w, r, &context)
	case r.Method == http.MethodPost && page.Post != nil:
		return page.handlePost(w, r, &context)
	default:
		return nil, ErrMethodNotAllowed
	}

	return context.action, nil
}

