package web

import "strings"

type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

/*func (r *router) findMdls(root *node, segs []string) []Middleware {
	queue := []*node{root}
	res := make([]Middleware, 0, 10)
	for i := 0; i < len(segs); i++ {
		seg := segs[i]
		var children []*node
		for _, cur := range queue {
			if len(cur.mdls) > 0 {
				res = append(res, cur.mdls...)
			}
			//children = append(children, cur.childOf(seg)...)
		}
		queue = children

	}

	for _, cur := range queue {
		if len(cur.mdls) > 0 {
			res = append(res, cur.mdls...)
		}
	}
	return res
}*/

func (r *router) AddRoute(method string, path string, handleFunc HandleFunc) {

	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	if path[0] != '/' {
		panic("web：路径必须以/ 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web：路径不能 以/ 结尾")
	}

	if path == "/" {
		if root.handler != nil {
			panic("重复注册 路由")
		}
		root.handler = handleFunc
		return
	}
	path = path[1:]
	segs := strings.Split(path, "/")
	for _, seg := range segs {

		if seg == "" {
			panic("web 不能有连续的//")
		}

		children := root.childOrCreate(seg)
		root = children
	}
	if root.handler != nil {
		panic("web 路由冲突 重复注册")
	}
	root.handler = handleFunc

}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{n: root}, true
	}
	//pathParams := make(map[string]string, 1)
	var pathParams map[string]string
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string, 1)
			}
			pathParams[child.path[1:]] = seg
		}
		root = child
	}
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true

}

type node struct {
	path       string
	children   map[string]*node
	paramChild *node
	starChild  *node
	handler    HandleFunc
	mdls       []Middleware
}

// 第二个bool 值是否是 param值匹配
// 第三个bool 是否命中
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil

	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, ok
}

func (n *node) childOrCreate(seg string) *node {
	if seg[0] == ':' {

		if n.starChild != nil {
			panic("web :不允许同时注册路径参数和通配符匹配， 已有通配符匹配")
		}

		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}
	if seg == "*" {

		if n.paramChild != nil {
			panic("web :不允许同时注册路径参数和通配符匹配， 已有参数匹配")
		}

		n.starChild = &node{
			path: seg,
		}
		return n.starChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
		/*res := &node{
			path: seg,
		}
		n.children[seg] = res
		return res*/
	}

	res, ok := n.children[seg]
	if !ok {
		res = &node{path: seg}
		n.children[seg] = res
	}
	return res
}

/*func (n *node) childOrCreate(seg string) interface{} {

}
*/

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
