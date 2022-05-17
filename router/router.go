package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nethttp-server/context"
)

/*
{
	"GET": {
		"pattern": http.HandlerFunc
		"pattern": http.HandlerFunc
	},
	"POST": {
		"pattern": http.HandlerFunc
		"pattern": http.HandlerFunc
	}
}
*/
type Router struct {
	Handlers map[string]map[string]context.HandlerFunc
}

func Match(pattern, path string) (bool, map[string]string) {
	if pattern == path {
		return true, nil
	}

	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	if len(patterns) != len(paths) {
		return false, nil
	}

	params := make(map[string]string)

	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
			continue
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			params[patterns[i][1:]] = paths[i]
		default:
			return false, nil
		}
	}

	return true, params

}


func (r *Router) HandleFunc(method, pattern string, h context.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	m, ok := r.Handlers[method]
	if !ok {
		m = make(map[string]context.HandlerFunc)
		r.Handlers[method] = m
	}
	m[pattern] = h
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for pattern, handler := range r.Handlers[req.Method] {
		if ok, params := Match(pattern, req.URL.Path); ok {
			c := context.Context{
				Params: make(map[string]interface{}),
				ResponseWriter: w,
				Request: req,
			}

			for k, v := range params {
				c.Params[k] = v
			}

			handler(&c)
			return
		}
	}
	http.NotFound(w, req)
	return
}
