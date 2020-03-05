package pages

import (
	"sweb/ccf"
	"net/url"
)

type Context struct {
	CommonContext

	GetParameters  url.Values
	PostParameters url.Values

	SessionParameters *ccf.Cookie
}
