package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
		htmlTemplates *template.Template // for html render
		funcMap       template.FuncMap   // for html render
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

// createStaticHandler creates static handler for Static()
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static serves static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// SetFuncMap renders function for custom
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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
	ctx.engine = engine
	engine.router.handle(ctx)
}