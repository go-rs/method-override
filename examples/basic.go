package main

import (
	"log"
	"net/http"

	methodoverride "github.com/go-rs/method-override"
	"github.com/go-rs/rest-api-framework"
)

func main() {
	var api = rest.New("/")
	api.Use(methodoverride.Load())

	api.Get("/", func(ctx *rest.Context) {
		ctx.JSON(ctx.Query)
	})

	log.Fatal(http.ListenAndServe(":8080", api))
}
