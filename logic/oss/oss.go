package oss

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hailongz/kk-lib/dynamic"
	"github.com/hailongz/kk-lib/http"
	"github.com/hailongz/kk-logic/logic"
)

type OSSGetLogic struct {
	logic.Logic
}

func (L *OSSGetLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	stype := dynamic.StringValue(L.Get(ctx, app, "type"), "string")
	path := dynamic.StringValue(L.Get(ctx, app, "path"), "")

	endpoint := dynamic.StringValue(L.Get(ctx, app, "endpoint"), "")
	accessKey := dynamic.StringValue(L.Get(ctx, app, "accessKey"), "")
	secretKey := dynamic.StringValue(L.Get(ctx, app, "secretKey"), "")
	bucket := dynamic.StringValue(L.Get(ctx, app, "bucket"), "")

	client, err := oss.New(endpoint, accessKey, secretKey)

	if err != nil {
		return L.Error(ctx, app, err)
	}

	buk, err := client.Bucket(bucket)

	if err != nil {
		return L.Error(ctx, app, err)
	}

	rd, err := buk.GetObject(path)

	if err != nil {
		return L.Error(ctx, app, err)
	}

	defer rd.Close()

	b, err := ioutil.ReadAll(rd)

	if err != nil {
		return L.Error(ctx, app, err)
	}

	switch stype {
	case "byte":
		ctx.Set(logic.ResultKeys, b)
	default:
		ctx.Set(logic.ResultKeys, string(b))
	}

	return L.Done(ctx, app, "done")
}

type OSSPutLogic struct {
	logic.Logic
}

func (L *OSSPutLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	content := L.Get(ctx, app, "content")
	file := L.Get(ctx, app, "file")
	url := L.Get(ctx, app, "url")

	if content != nil || file != nil || url != nil {

		path := dynamic.StringValue(L.Get(ctx, app, "path"), "")

		endpoint := dynamic.StringValue(L.Get(ctx, app, "endpoint"), "")
		accessKey := dynamic.StringValue(L.Get(ctx, app, "accessKey"), "")
		secretKey := dynamic.StringValue(L.Get(ctx, app, "secretKey"), "")
		bucket := dynamic.StringValue(L.Get(ctx, app, "bucket"), "")

		client, err := oss.New(endpoint, accessKey, secretKey)

		if err != nil {
			return L.Error(ctx, app, err)
		}

		buk, err := client.Bucket(bucket)

		if err != nil {
			return L.Error(ctx, app, err)
		}

		if file != nil {

			fd, ok := file.(*multipart.FileHeader)

			if ok {

				rd, err := fd.Open()

				if err != nil {
					return L.Error(ctx, app, err)
				}

				err = buk.PutObject(path, rd)

				rd.Close()

				log.Println("[OSS]", path, err)

				if err != nil {
					return L.Error(ctx, app, err)
				}
			} else {
				return L.Error(ctx, app, logic.NewError(logic.ERROR_UNKNOWN, "未找到文件内容"))
			}
		} else if url != nil {

			u := dynamic.StringValue(url, "")

			if strings.HasPrefix(u, "http") {

				options := http.Options{}

				options.Url = u
				options.ResponseType = http.OptionResponseTypeByte
				options.Method = "GET"

				b, err := http.Send(&options)

				if err != nil {
					return L.Error(ctx, app, err)
				}

				err = buk.PutObject(path, bytes.NewReader(b.([]byte)))

				log.Println("[OSS]", url, path, err, len(b.([]byte)))

				if err != nil {
					return L.Error(ctx, app, err)
				}

			} else {
				return L.Error(ctx, app, logic.NewError(logic.ERROR_UNKNOWN, "不支持的URL"))
			}

		} else if content != nil {

			b, ok := content.([]byte)

			if !ok {

				b = []byte(dynamic.StringValue(content, ""))
			}

			err = buk.PutObject(path, bytes.NewReader(b))

			log.Println("[OSS]", url, path, err, len(b))

			if err != nil {
				return L.Error(ctx, app, err)
			}

		}

	} else {
		return L.Error(ctx, app, logic.NewError(logic.ERROR_UNKNOWN, "未找到内容"))
	}

	return L.Done(ctx, app, "done")
}

type OSSURLLogic struct {
	logic.Logic
}

func (L *OSSURLLogic) Exec(ctx logic.IContext, app logic.IApp) error {

	L.Logic.Exec(ctx, app)

	method := dynamic.StringValue(L.Get(ctx, app, "method"), "GET")
	path := dynamic.StringValue(L.Get(ctx, app, "path"), "")
	expires := dynamic.IntValue(L.Get(ctx, app, "expires"), 300)
	contentType := dynamic.StringValue(L.Get(ctx, app, "contentType"), "")

	endpoint := dynamic.StringValue(L.Get(ctx, app, "endpoint"), "")
	accessKey := dynamic.StringValue(L.Get(ctx, app, "accessKey"), "")
	secretKey := dynamic.StringValue(L.Get(ctx, app, "secretKey"), "")
	bucket := dynamic.StringValue(L.Get(ctx, app, "bucket"), "")

	client, err := oss.New(endpoint, accessKey, secretKey)

	if err != nil {
		return L.Error(ctx, app, err)
	}

	buk, err := client.Bucket(bucket)

	if err != nil {
		return L.Error(ctx, app, err)
	}

	m := oss.HTTPGet

	switch method {
	case "POST":
		m = oss.HTTPPost
		break
	case "PUT":
		m = oss.HTTPPut
		break
	case "DELETE":
		m = oss.HTTPDelete
		break
	}

	options := []oss.Option{}

	if contentType != "" {
		options = append(options, oss.ContentType(contentType))
	}

	u, err := buk.SignURL(path, m, expires, options...)

	if err != nil {
		return L.Error(ctx, app, err)
	}

	ctx.Set(logic.ResultKeys, u)

	return L.Done(ctx, app, "done")
}

func init() {
	logic.Openlib("kk.Logic.OSS.Get", func(object interface{}) logic.ILogic {
		v := OSSGetLogic{}
		v.Init(object)
		return &v
	})
	logic.Openlib("kk.Logic.OSS.Put", func(object interface{}) logic.ILogic {
		v := OSSPutLogic{}
		v.Init(object)
		return &v
	})
	logic.Openlib("kk.Logic.OSS.URL", func(object interface{}) logic.ILogic {
		v := OSSURLLogic{}
		v.Init(object)
		return &v
	})
}
