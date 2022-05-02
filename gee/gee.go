package gee

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(ctx *Context)

// Engine implements the Handler interface.
// it uses router to store handlers.
// Engine is the sub-class of RouterGroup
type (
	RouterGroup struct {
		prefix string
		middlewares []HandlerFunc // support middleware
		engine *Engine			 // all groups share one engine instance
	}

	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup	// store all groups
	}
)

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	group := &RouterGroup{engine: engine}
	engine.RouterGroup = group
	engine.groups = []*RouterGroup{group}
	return engine
}

// Group is the constructor of RouterGroup
// all groups share an engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup{
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// addRoute adds handler function into Engine.router.handlers,
// invoked by methods like GET(), POST()
func (group *RouterGroup) addRoute(method string, suffix string, handler HandlerFunc) {
	pattern := group.prefix + suffix
	log.Printf("Route %4s - %s\n", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request handler
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request handler
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP defines how engine route requests
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		// Engine's prefix is ""(default string), and every string has this prefix,
		// so engine's middleware(has no other groups) is global
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	ctx := newContext(w, req)
	ctx.handlers = middlewares
	engine.router.handle(ctx)
}