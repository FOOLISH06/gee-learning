package gee

import (
	"log"
	"net/http"
)

// router is a struct contains handler functions and handle them
// router and its methods are not public. They are only invoked by Engine.
// key of handlers is combination of request method and request path, like "GET-/hello"
type router struct {
	handlers map[string]HandlerFunc
}

// newRouter returns a new router for Engine
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// addRoute adds handler function into router,
// uses request method and path as key
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle is the real method that calls handler functions.
// Engine.ServeHTTP creates a Context with http.ResponseWriter and http.Request,
// and then calls router.handle to handle the request
func (r *router) handle(ctx *Context) {
	key := ctx.Method + "-" + ctx.Path
	if handler, ok := r.handlers[key]; ok {
		handler(ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
		log.Printf("handler %q not found\n", key)
	}
}