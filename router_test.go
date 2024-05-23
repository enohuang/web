package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(r)
	/*_ = msg
	_ = ok
	if !ok {
		t.Errorf(msg)
	}*/
	assert.True(t, ok, msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.AddRoute(http.MethodGet, "", mockHandler)
	}, "web:路径必须以/开头")

	assert.Panicsf(t, func() {
		r.AddRoute(http.MethodGet, "/a/b/c/", mockHandler)
	}, "web:路径不能以/结尾")

	assert.Panicsf(t, func() {
		r.AddRoute(http.MethodGet, "/a//b/c/", mockHandler)
	}, "web:路径不能有//")

	r = newRouter()
	r.AddRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.AddRoute(http.MethodGet, "/", mockHandler)
	}, "重复注册")

	r = newRouter()
	r.AddRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panicsf(t, func() {
		r.AddRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "重复注册/a/b/c")

}

func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的 http method"), false
		}
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不相等"), false
	}

	//比较handler
	hHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if yHandler != hHandler {
		return fmt.Sprintf("Handler 不相等"), false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点%s 不存在", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	testRoute := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
	}

	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRoute {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		wantNode  *node
	}{
		{
			name:      "user home",
			method:    http.MethodGet,
			path:      "/user/home",
			wantFound: true,
			wantNode: &node{
				/*handler: mockHandler,*/
				path:    "home",
				handler: mockHandler,
			},
		},
		{
			name:      "method not found",
			method:    http.MethodOptions,
			path:      "/order/detail",
			wantFound: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			msg, ok := tc.wantNode.equal(n.n)
			assert.True(t, ok, msg)

			//assert.Equal(t, tc.wantNode.path, n.path)
			//assert.Equal(t, tc.wantNode.children, n.children)

			/*nHandler := reflect.ValueOf(n.handler)
			yHandler := reflect.ValueOf(n.handler)
			assert.True(t, nHandler == yHandler)*/

		})
	}
}
