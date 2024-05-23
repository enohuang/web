package web

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileUpLoader struct {
	FileField   string
	DstPathFunc func(*multipart.FileHeader) string
}

func (u FileUpLoader) Handle() HandleFunc {
	return func(ctx *Context) {
		//上传文件逻辑

		//第一步：读到文件内容
		// 目标路径
		// 保存文件
		// 返回响应

		file, fileHeader, err := ctx.Req.FormFile(u.FileField)
		if err != nil {
			ctx.StatusCode = 500
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer file.Close()

		dst := u.DstPathFunc(fileHeader)
		fmt.Println("创建路径", dst)
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			ctx.StatusCode = 500
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer dstFile.Close()
		io.CopyBuffer(dstFile, file, nil)
		ctx.StatusCode = 200
		ctx.RespData = []byte("上传成功")

	}
}

type FileDownloader struct {
	Dir string
}

func (d *FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespData = []byte("找不到目标文件")
			ctx.StatusCode = http.StatusBadRequest
			return
		}

		dst := filepath.Join(d.Dir, req)
		path := filepath.Join(d.Dir, filepath.Clean(req))

		//做一个校验 防止 用户上传 ../../../file.txt 逃逸testdata/download之外
		dst, err = filepath.Abs(dst)
		if strings.Contains(dst, d.Dir) {
			//校验是否逃逸
		}

		//fmt.Println(dst)
		fn := filepath.Base(dst)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, path)

	}
}

type StaticResourceHandlerOption func(handler *StaticResourceHandler)

func StaticWithMaxSize(maxSize int) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.maxSize = maxSize
	}
}

type StaticResourceHandler struct {
	dir               string
	cache             *lru.Cache
	extContextTypeMap map[string]string
	//大文件不缓存
	maxSize int
}

func NewStaticResourceHandler(dir string, opts ...StaticResourceHandlerOption) *StaticResourceHandler {

	c, _ := lru.New(1024 * 1024)
	res := &StaticResourceHandler{
		dir:   dir,
		cache: c,
		extContextTypeMap: map[string]string{
			"jpg": "image/jpeg",
			"png": "image/png",
			"pdf": "image/pdf",
			"jpe": "image/jpeg",
		},
	}

	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (s *StaticResourceHandler) HandleWithCache(ctx *Context) {
	//1. 拿到目标文件名

	//2.定位到目标文件，并且读出来

	//3.返回前端

	file, err := ctx.PathValue("file")
	if err != nil {
		ctx.StatusCode = http.StatusBadRequest
		ctx.RespData = []byte("请求路径不对")
		return
	}

	dst := filepath.Join(s.dir, file)
	ext := filepath.Ext(dst)
	header := ctx.Resp.Header()

	if data, ok := s.cache.Get(file); ok {
		header.Set("Content-Type", s.extContextTypeMap[ext[1:]])
		header.Set("Content-Length", strconv.Itoa(len(data.([]byte))))
		// 可能是文本文件， 图片，多媒体

		ctx.RespData = data.([]byte)
		ctx.StatusCode = 200
		return
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		ctx.StatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器错误")
	}

	fmt.Println("dst", dst, "data", len(data))
	//fmt.Println(data)

}

func (s *StaticResourceHandler) Handle(ctx *Context) {
	//1. 拿到目标文件名

	//2.定位到目标文件，并且读出来

	//3.返回前端

	file, err := ctx.PathValue("file")
	if err != nil {
		ctx.StatusCode = http.StatusBadRequest
		ctx.RespData = []byte("请求路径不对")
		return
	}

	dst := filepath.Join(s.dir, file)

	data, err := os.ReadFile(dst)
	if err != nil {
		ctx.StatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器错误")
	}

	fmt.Println("dst", dst, "data", len(data))
	//fmt.Println(data)
	ext := filepath.Ext(dst)
	header := ctx.Resp.Header()
	header.Set("Content-Type", s.extContextTypeMap[ext[1:]])
	header.Set("Content-Length", strconv.Itoa(len(data)))
	// 可能是文本文件， 图片，多媒体

	ctx.RespData = data
	ctx.StatusCode = 200

}
