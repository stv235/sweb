package pages

import (
	"sweb/ccf"
	"net/http"
	"net/url"
)

type Route interface {
	handle(w http.ResponseWriter, r *http.Request, pathParameters url.Values, sessionParameters *ccf.Cookie) (action, error)
}