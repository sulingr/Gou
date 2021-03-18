package main

import (
	"gou"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gou.HandlerFunc {
	return func(c *gou.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
func main() {
	r := gou.New()
	r.Static("/assets", "./static")
	r.Use(gou.Logger())
	r.GET("/index", func(c *gou.Context) {

	})
	v1 := r.Team("/v1")
	{
		v1.GET("/", func(c *gou.Context) {

		})

		v1.GET("/hello", func(c *gou.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Team("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gou.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
