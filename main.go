package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nethttp-server/context"
	"github.com/nethttp-server/middleware"
	"github.com/nethttp-server/router"
)

func main() {
	emptyMap := make(map[string]map[string]context.HandlerFunc)

	r := &router.Router{Handlers: emptyMap}

	r.HandleFunc("GET", "/", middleware.LogHandler(func(c *context.Context) {
		fmt.Fprintln(c.ResponseWriter, "GET /")
	}))

	r.HandleFunc("GET", "/about", func(c *context.Context) {
		fmt.Fprintln(c.ResponseWriter, "GET about")
	})
	r.HandleFunc("GET", "/users/:id", middleware.LogHandler(middleware.RecoverHandler(
		func(c *context.Context) {
			if c.Params["id"] == "0" {
				panic("id is zero")
			}
			fmt.Fprintln(c.ResponseWriter, "retrieve user", c.Params["id"])
		},
	)))
	r.HandleFunc("GET", "/users/:id/addresses/:addresses_id", func(c *context.Context) {
		fmt.Fprintln(c.ResponseWriter, "retrieve user", c.Params["id"], "address", c.Params["addresses_id"])
	})
	r.HandleFunc("POST", "/users",
		middleware.LogHandler(
			middleware.RecoverHandler(
				middleware.ParseFormHandler(
					middleware.ParseJsonBodyHandler(
						func(c *context.Context) {
							fmt.Fprintln(c.ResponseWriter, "create user")
						},
					),
				),
			),
		),
	)

	r.HandleFunc("POST", "/users/:id/addresses", func(c *context.Context) {
		fmt.Fprintln(c.ResponseWriter, "create user's address", c.Params["id"])
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
