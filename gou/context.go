package gou

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	//原始对象
	ReW http.ResponseWriter
	Req *http.Request
	//请求信息
	Path   string
	Method string
	//响应信息
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		ReW:    w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.ReW.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.ReW.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.ReW.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.ReW)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.ReW, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.ReW.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.ReW.Write([]byte(html))
}
