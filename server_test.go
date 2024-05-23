//go:build e2e

package web

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
)

func TestServer(t *testing.T) {
	//
	h := NewHTTPServer() //&HTTPServer{}
	h.AddRoute(http.MethodGet, "/user", func(ctx *Context) {
		fmt.Println("处理第一件事")
		fmt.Println("处理第二件事")
	})

	handler1 := func(ctx *Context) {
		fmt.Println("处理第一件事")
	}

	handler2 := func(ctx *Context) {
		fmt.Println("处理第二件事")
	}

	h.AddRoute(http.MethodGet, "/user2", func(ctx *Context) {
		handler1(ctx)
		handler2(ctx)
	})

	h.AddRoute(http.MethodGet, "/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello order detail"))
	})

	h.AddRoute(http.MethodPost, "/from", func(ctx *Context) {
		ctx.Req.ParseForm()
		ctx.Resp.Write([]byte("hello order detail"))
	})

	h.Start(":8081")

}

type SafeContext struct {
	Context
	mutex sync.RWMutex
}

func (c *SafeContext) RespJSONOK(val any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.RespJSON(http.StatusOK, val)
}

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHTTPServer()
	server.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个 before")
				next(ctx)
				fmt.Println("第一个 after")
			}
		},

		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个 before")
				next(ctx)
				fmt.Println("第二个 after")
			}
		},

		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个 中断")

			}
		},

		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四个个 中断")

			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})

}
