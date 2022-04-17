package gee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Context is a struct contains information of request and response,
// created by Engine and passed to handler function when ServeHTTP().
// Context is created when a request is coming and dropped when response is sent
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req *http.Request
	// request info
	Path string
	Method string
	Params map[string]string
	// response info
	StatusCode int
}

// H is an alias, used to  write json data to response
type H map[string]interface{}

// newContext is not public, invoked by Engine when ServeHTTP()
func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: writer,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
	}
}

// Status writes status code into header,
// after invoke Status(), can't modify header anymore,
// so set other headers like "content-type" before
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader sets response header with key-value
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// Query gets request query parameters by key
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// PostForm gets request form data by key
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Param gets the parsed route param
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// String writes a string with format into response,
// values is a variable-length argument list, type of interface{}
// values... represent like a slice
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("content-type", "text/plain")
	c.Status(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		// can't use http.error(), otherwise it will call WriteHeader() twice
		log.Printf("String() response string failed, err: %q\n", err.Error())
	}
}

// JSON writes a json data into response,
// obj is usually type of gee.H
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("content-type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		log.Printf("JSON() encode failed, err: %q, data: %v\n", err.Error(), obj)
	}
}

// Data writes []byte data into response
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, err := c.Writer.Write(data)
	if err != nil {
		log.Printf("Data() response data failed, err: %q\n", err.Error())
	}
}

// HTML writes html statement into response
func (c *Context) HTML(code int, html string) {
	c.SetHeader("content-type", "text/html")
	c.Status(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		log.Printf("HTML() response data failed, err: %q\n", err.Error())
	}
}
