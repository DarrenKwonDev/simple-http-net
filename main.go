package main

import (
	"fmt"
	"net/http"

	"github.com/nethttp-server/router"
)

type A struct {
	name string
}

func main() {
	emptyMap := make(map[string]map[string]http.HandlerFunc)

	r := &router.Router{Handlers: emptyMap}

	r.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GET /")
	})
	r.HandleFunc("GET", "/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GET about")
	})
	r.HandleFunc("GET", "/users/:id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "retrieve user")
	})
	r.HandleFunc("GET", "/users/:id/addresses/:addresses_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "retrieve user")
	})
	r.HandleFunc("POST", "/users", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "create user")
	})
	r.HandleFunc("POST", "/users/:id/addresses", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "create user's address")
	})

	fmt.Println("Server is running at http://localhost:8080/")
	http.ListenAndServe(":8080", r)
}
