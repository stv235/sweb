package pages

import (
	"fmt"
	"log"
	"net/url"
)

type CommonContext struct {
	action action

	Path           string
	PathParameters url.Values
}

func (ctx *CommonContext) Render(name string, data interface{}) {
	ctx.checkAction()
	ctx.action = &renderer{name: name, data:data}
}

func (ctx *CommonContext) SendFile(buf []byte, mimeType string) {
	ctx.checkAction()

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	ctx.action = &file{ buf: buf, mimeType: mimeType }
}

func (ctx *CommonContext) Redirect(url string, parameters ...interface{}) {
	ctx.checkAction()

	if len(parameters) > 0 {
		url = fmt.Sprintf(url, parameters...)
	}

	ctx.action = &redirector{ url: url }
}

func (ctx *CommonContext) Status(statusCode int, message string) {
	ctx.checkAction()
	ctx.action = &status{ code: statusCode, message: message }
}

func (ctx *CommonContext) checkAction() {
	if ctx.action != nil {
		log.Fatalln("page action already set, did you call a page action twice?")
	}
}

