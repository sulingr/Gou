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
	Params map[string]string
	//响应信息
	StatusCode int
	//中间件
	handlers []HandlerFunc
	index    int //记录当前执行到第几个中间件
	//engine 指针
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		ReW:    w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
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

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, map[string]interface{}{"message": err})
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

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.ReW, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
