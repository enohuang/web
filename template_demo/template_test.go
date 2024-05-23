package template_demo

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	type User struct {
		Name string
	}
	tpl := template.New("hello-world")
	//. 当前作用域的当前对象
	tpl, err := tpl.Parse(`Hello, {{.Name}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, User{Name: "Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())
}

func TestMap(t *testing.T) {
	type User struct {
		Name string
	}
	tpl := template.New("hello-world")
	//. 当前作用域的当前对象
	tpl, err := tpl.Parse(`Hello, {{.Name}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, map[string]string{"Name": "Tom"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, Tom`, buffer.String())

}

func TestService(t *testing.T) {
	type User struct {
		Name string
	}
	tpl := template.New("hello-world")
	//. 当前作用域的当前对象
	tpl, err := tpl.Parse(`
 {{- $service := .GenName -}}
type {{ $service }} struct{
	Path string 
}
    `)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, map[string]string{"GenName": "Tom"})
	require.NoError(t, err)
	t.Log(buffer.String())
	//assert.Equal(t, `Hello, Tom`, buffer.String())
}

func TestFuncCall(t *testing.T) {
	tpl := template.New("hello-world")
	//. 当前作用域的当前对象
	tpl, err := tpl.Parse(`
切片长度  {{len .Slice}}
{{printf "%.2f"  1.2345}}
Hello, {{.Hello "Tom" "Jerry" }}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{Slice: []string{"a", "b"}})
	require.NoError(t, err)
	assert.Equal(t, `
切片长度  2
1.23
Hello, Tom.Jerry`, buffer.String())
}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(first string, last string) string {
	return fmt.Sprintf("%s.% s", first, last)
}

func TestForLoop(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $ele  := .Slice}}
{{- .}}
{{$idx}}-{{$ele}}
{{end}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	t.Log(buffer.String())
}

func TestForLoop2(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $ele  := .}}
{{- $idx}}
{{- end}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, make([]int, 100))
	t.Log(buffer.String())
}

func TestIfSlse(t *testing.T) {
	type User struct {
		Age int
	}
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (lt .Age 6)}}
儿童：(0, 6]
{{ else if and (gt .Age 6) (lt .Age 18)}}
少年: (6,18]
{{ else }}
成人:  >18
{{end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, User{Age: 19})
	require.NoError(t, err)
	assert.Equal(t, ``, buffer.String())
}
