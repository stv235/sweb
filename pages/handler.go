package pages

import (
	"fmt"
	"sweb/ccf"
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Handler struct {
	sessionConfig *SessionConfig
	routeHandlers []RouteHandler
	html *template.Template

	ErrorFn       func(ctx *ErrorContext)
}

func (handler Handler) Match(r *http.Request) (Route, url.Values) {
	// TODO: pre-split path?
	for _, routeHandler := range handler.routeHandlers {
		if ok, parameters := routeHandler.Match(r.URL.Path); ok {
			return routeHandler.Route, parameters
		}
	}

	return nil, nil
}

func (handler Handler) serveAction(action action, w http.ResponseWriter, r *http.Request) {
	switch action.(type) {
	case *redirector:
		redirector, _ := action.(*redirector)
		http.Redirect(w, r, redirector.url, http.StatusFound)
	case *renderer:
		renderer, _ := action.(*renderer)

		buf := bytes.NewBuffer(nil)
		err := handler.html.ExecuteTemplate(buf, renderer.name + ".html", renderer.data)

		if err != nil {
			panic(err)
		}

		if _, err := w.Write(buf.Bytes()); err != nil {
			log.Println("[HTTP]", err)
		}
	case *file:
		file, _ := action.(*file)

		w.Header().Add("Content-Type", file.mimeType)

		_, err := w.Write(file.buf)

		if err != nil {
			log.Println("[HTTP]", err)
		}
	case *status:
		status, _ := action.(*status)

		w.WriteHeader(status.code)
		if _, err := w.Write([]byte(status.message)); err != nil {
			log.Println("[HTTP]", err)
		}
	case *send:
		if action.(*send).closeFn != nil {
			defer action.(*send).closeFn()
		}

		if _, err := io.Copy(w, action.(*send).r); err != nil {
			log.Println(TagHttp, err)
		}
	case *requireAuth:
		w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", action.(*requireAuth).realm))
		w.WriteHeader(401)

		if _, err := w.Write([]byte("Unauthorised.\n")); err != nil {
			log.Println(TagHttp, err)
		}

		return
	}
}

func (handler Handler) saveSession(session ccf.Cookie, w http.ResponseWriter, r* http.Request) {
	if handler.sessionConfig != nil {
		timeout := time.Now().Add(time.Duration(handler.sessionConfig.Timeout) * time.Minute)
		val, err := session.Encode(timeout, handler.sessionConfig.EncryptKey, handler.sessionConfig.SignKey)

		if err != nil {
			log.Panicln("[HTTP]", "cookie", err)
		} else {
			cookie := &http.Cookie{}
			cookie.Path = "/"
			cookie.HttpOnly = true
			cookie.Name = handler.sessionConfig.CookieName
			cookie.Value = val

			http.SetCookie(w, cookie)
		}
	}
}

func (handler Handler) deleteSession(w http.ResponseWriter) {
	if handler.sessionConfig != nil {
		cookie := &http.Cookie{}
		cookie.Name = handler.sessionConfig.CookieName
		cookie.Value = ""
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
	}
}

func (handler Handler) serveRoute(w http.ResponseWriter, r *http.Request, route Route, pathParameter url.Values) {
	defer func() {
		defer func() {
			if e := recover(); e != nil {
				switch e.(type) {
				case string:
					log.Println("[LOGIC]", e)
				case error:
					log.Println("[LOGIC]", e.(error).Error())
				}

				// Remove session cookie on error
				handler.deleteSession(w)
				w.WriteHeader(http.StatusInternalServerError)

				w.Write([]byte("Internal error"))
			}
		}()

		if err := recover(); err != nil {
			ctx := ErrorContext{}
			ctx.Err = err

			handler.ErrorFn(&ctx)
			handler.serveAction(ctx.action, w, r)
		}
	}()

	sessionParameters := ccf.Cookie{}
	handler.loadSession(r, &sessionParameters)

	action, err := route.handle(w, r, pathParameter, &sessionParameters)

	if err != nil {
		log.Panicln(TagHttp, err)
	}

	if action != nil {
		handler.saveSession(sessionParameters, w, r)
		handler.serveAction(action, w, r)
	} else {
		log.Println(TagLogic, "no action")
	}
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[HTTP]", r.Method, r.URL.Path)

	route, pathParameters := handler.Match(r)

	if route != nil {
		handler.serveRoute(w, r, route, pathParameters)
	} else {
		log.Println(TagHttp, "unknown page", r.URL)
		w.WriteHeader(http.StatusNotFound)

		if _, err := io.WriteString(w, "Unknown page"); err != nil {
			log.Panicln(TagLogic, err)
		}
	}
}

func (handler *Handler) appendRoute(pattern string, route Route) {
	routeHandler := RouteHandler{ Route: route }

	if pattern == "" {
		routeHandler.patternParts = make([]string, 0)
	} else {
		routeHandler.patternParts = strings.Split(pattern, "/")
	}

	handler.routeHandlers = append(handler.routeHandlers, routeHandler)
}

func (handler *Handler) RegisterPage(pattern string, page Page) {
	handler.appendRoute(pattern, page)
}

func (handler *Handler) RegisterStream(pattern string, stream Stream) {
	handler.appendRoute(pattern, stream)
}

func (handler Handler) loadSession(r *http.Request, sessionParameters *ccf.Cookie) {
	if handler.sessionConfig != nil {
		if cookie, err := r.Cookie(handler.sessionConfig.CookieName); err == nil && cookie != nil {
			if err := sessionParameters.Decode(cookie.Value, handler.sessionConfig.EncryptKey, handler.sessionConfig.SignKey); err != nil {
				// invalid cookie is possible - discard
				log.Println(TagHttp, err)
			}
		}
	}
}

func NewHandler(html *template.Template, sessionConfig *SessionConfig, errorFn func(ctx *ErrorContext)) *Handler {
	if sessionConfig != nil {
		if sessionConfig.CookieName == "" {
			sessionConfig.CookieName = "pages_auth"
		}

		if sessionConfig.EncryptKey == "" {
			log.Fatalln("[HTTP]", "no session encryption key")
		}

		if sessionConfig.SignKey == "" {
			log.Fatalln("[HTTP]", "no session sign key")
		}

		if sessionConfig.Timeout == 0 {
			log.Fatalln("[HTTP]", "invalid session timeout")
		}
	}

	handler := Handler{}
	handler.ErrorFn = errorFn
	handler.html = html
	handler.sessionConfig = sessionConfig

	return &handler
}