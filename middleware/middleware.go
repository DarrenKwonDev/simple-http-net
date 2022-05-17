package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nethttp-server/context"
)

type Middleware func(next context.HandlerFunc) context.HandlerFunc

func LogHandler(next context.HandlerFunc) context.HandlerFunc {
	return func(c *context.Context) {
		t := time.Now()
		next(c)
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Since(t))
	}
}

func RecoverHandler(next context.HandlerFunc) context.HandlerFunc {
	return func(c *context.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(c.ResponseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next(c)
	}
}

func ParseFormHandler(next context.HandlerFunc) context.HandlerFunc {
	return func(c *context.Context) {
		c.Request.ParseForm()
		fmt.Println(c.Request.PostForm)
		for k, v := range c.Request.PostForm {
			if len(v) > 0 {
				c.Params[k] = v[0]
			}
		}
		next(c)
	}
}

func ParseJsonBodyHandler(next context.HandlerFunc) context.HandlerFunc {
	return func(c *context.Context) {
		var m map[string]interface{}
		if json.NewDecoder(c.Request.Body).Decode(&m); len(m) > 0 {
			for k, v := range m {
				c.Params[k] = v
			}
		}
		next(c)
	}
}