package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
)

func main() {
	router := gee.New()
	router.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "requset URL path: %q", req.URL.Path)
	})
	router.GET("/headers", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	err := router.Run(":8080")
	if err != nil {
		log.Println("engine run failed, error: ", err.Error())
	}
}