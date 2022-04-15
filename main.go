package main

import (
	"gee"
	"log"
	"net/http"
)

func main() {
	router := gee.New()

	router.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>hello world<h1>")
	})
	router.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %q\n", c.Query("name"), c.Path)
	})
	router.GET("/query", func(ctx *gee.Context) {
		ctx.JSON(http.StatusOK, gee.H{
			"name": ctx.Query("name"),
			"age": ctx.Query("age"),
		})
	})
	router.POST("/user", func(ctx *gee.Context) {
		ctx.JSON(http.StatusOK, gee.H{
			"name": ctx.PostForm("name"),
			"age": ctx.PostForm("age"),
		})
	})

	err := router.Run(":8080")
	if err != nil {
		log.Println("engine run failed, error: ", err.Error())
	}
}