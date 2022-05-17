package middleware

import (
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