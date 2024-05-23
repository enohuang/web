package web

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"testing"
)

func TestUpload2(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)

	engine := &GoTemplateEngine{
		T: tpl,
	}

	h := NewHTTPServer(ServerWithTemplateEngine(engine))
	h.AddRoute(http.MethodGet, "/upload", func(ctx *Context) {
		err := ctx.Render("upload.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
		//ctx.Resp.Write([]byte("chenggong"))
	})

	fu := FileUpLoader{
		FileField: "myfile",
		DstPathFunc: func(header *multipart.FileHeader) string {
			return filepath.Join("testdata", "upload", header.Filename)
		},
	}

	h.AddRoute(http.MethodPost, "/upload", fu.Handle())

	h.Start(":8081")
}

func TestDownload(t *testing.T) {

	h := NewHTTPServer()
	fu := FileDownloader{
		Dir: filepath.Join("testdata", "download"), //考虑各个平台的兼容性
	}

	h.AddRoute(http.MethodGet, "/download", fu.Handle())

	h.Start(":8081")
}

func TestStaticResourceHandler_Handle(t *testing.T) {

	h := NewHTTPServer()
	fu := /*StaticResourceHandler{
			dir: filepath.Join("testdata", "static"),
		}*/
		NewStaticResourceHandler(filepath.Join("testdata", "static"))

	h.AddRoute(http.MethodGet, "/static/:file", fu.HandleWithCache)

	h.Start(":8081")
}
