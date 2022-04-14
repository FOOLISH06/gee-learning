package gee

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(w http.ResponseWriter, req *http.Request)

// Engine implements the Handler interface,
// use map to store handler functions,
// key is combination of request method and request path, like "GET-/hello"
type Engine struct {
	router map[string]HandlerFunc
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: make(map[string]HandlerFunc)}
	return engine
}

// addRoute adds handler function into engine.router,
// invoked by methods like GET(), POST()
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP defines how engine route requests
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(404)
		_, err := fmt.Fprintf(w, "404 NOT FOUNF: %s\n", req.URL)
		if err != nil {
			log.Println("write to ResponseWriter failed when 404, error: ", err)
			return
		}
	}
}