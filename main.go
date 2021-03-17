package main

import (
	"gou"
	"net/http"
)

func main() {
	r := gou.New()
	r.GET("/index", func(c *gou.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Team("/v1")
	{
		v1.GET("/", func(c *gou.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gou</h1>")
		})

		v1.GET("/hello", func(c *gou.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Team("/v2")
	{
		v2.GET("/hello/:name", func(c *gou.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gou.Context) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}
