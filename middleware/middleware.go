package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
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

func StaticHandler(next context.HandlerFunc) context.HandlerFunc {
	var(
		dir = http.Dir(".")
		indexFile = "index.html"
	)

	return func(c *context.Context) {
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			next(c)
			return
		}
		file := c.Request.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			next(c)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			next(c)
			return
		}

		// 디렉토리인 경우 indexFile을 사용하도록 처리
		if fi.IsDir() {
			if !strings.HasSuffix(c.Request.URL.Path, "/") {
				http.Redirect(c.ResponseWriter, c.Request, c.Request.URL.Path+"/", http.StatusFound)
				return
			}
		}

		file = path.Join(file, indexFile)
		
		f, err = dir.Open(file)
		if err != nil {
			next(c)
			return
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			next(c)
			return
		}

		// next 호출하지 않고 그냥 파일 서빙하고 끝냄.
		http.ServeContent(c.ResponseWriter, c.Request, file, fi.ModTime(), f)
	}
}