package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Context struct {
	Req          *http.Request
	Resp         http.ResponseWriter
	PathParams   map[string]string
	StatusCode   int
	RespData     []byte
	queryValue   url.Values
	MatchedRoute string
	tlpEngine    TemplateEngine
	UserValues   map[string]any
}

func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("web 输入不能为nil")
	}

	if c.Req.Body == nil {
		return errors.New("c.Req.Body == nil ")
	}
	decoder := json.NewDecoder(c.Req.Body)
	decoder.UseNumber()
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	VALS, OK := c.Req.Form[key]
	if OK {
		return "", errors.New("key not exist")
	}
	return VALS[0], nil
}

func (c *Context) QueryValue1(key string) StringValue {
	if c.queryValue == nil {
		c.queryValue = c.Req.URL.Query()
	}
	vals, ok := c.queryValue[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key 不存在"),
		}
	}
	return StringValue{val: vals[0]}
}

func (c *Context) QueryValue(key string) (string, error) {
	if c.queryValue == nil {
		c.queryValue = c.Req.URL.Query()
	}
	vals, ok := c.queryValue[key]
	if !ok {
		return "", errors.New("web: key 不存在")
	}
	return vals[0], nil
}

func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New(" web : key bucunzai ")
	}
	return val, nil
}

func (c *Context) PathValue1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New(" web : key bucunzai ")}
	}
	return StringValue{val: val}
}

type StringValue struct {
	val string
	err error
}

func (c *Context) RespJSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(status)
	c.Resp.Header().Set("Content-Type", "application/json")
	_, err = c.Resp.Write(data)
	return err
}

func (c *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(c.Resp, ck)
}

func (c *Context) Render(tplName string, data any) error {
	fmt.Println("Render................................", c.tlpEngine)

	var err error
	c.RespData, err = c.tlpEngine.Render(c.Req.Context(), tplName, data)
	if err != nil {
		c.StatusCode = 500
		fmt.Println("errr................................", err)
		return err
	}

	c.StatusCode = http.StatusOK
	return nil

}

// 这个泛型不行，因为创建的时候我们不知道用户需要什么 T
type StringValue2[T any] struct {
	val string
	err error
	a   T
}
