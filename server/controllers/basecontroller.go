package controllers

import (
	"bytes"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/web"
	"html/template"
	"net/http"
)

type BaseController struct {
	Controller
}

type baseControllerResult struct {
	Code int
	Data interface{}
	Msg  string
}

func (this *BaseController) jsonRender(result baseControllerResult) {
	resultString, err := EncodeJson(result)
	if err != nil {
		panic(err)
	}
	this.Ctx.Write(resultString)
}

func (this *BaseController) AutoRender(returnValue interface{}, viewname string) {
	result := baseControllerResult{}
	resultError, ok := returnValue.(Exception)
	request := this.Ctx.GetRawRequest().(*http.Request)
	responseWriter := this.Ctx.GetRawResponseWriter().(http.ResponseWriter)

	//带错误码的error
	result.Code = resultError.GetCode()
	result.Msg = resultError.GetMessage()
	if ok {
		//用户未登陆
		if result.Code == 10001 {
			this.Log.Debug("not signin")
			http.Redirect(responseWriter, request, "/index/signin", http.StatusMovedPermanently)
			return
		}

		this.Ctx.Write([]byte(resultError.GetMessage()))
	} else {
		//如果是跳转页面
		if viewname == "redirect" {
			http.Redirect(responseWriter, request, returnValue.(string), http.StatusMovedPermanently)
			return
		}

		//正常返回
		buffer := bytes.NewBuffer(nil)
		t, err := template.ParseFiles("../static/" + viewname + ".html")
		if err != nil {
			panic(err)
		}

		err = t.Execute(buffer, returnValue)
		if err != nil {
			panic(err)
		}
		this.Ctx.Write(buffer.Bytes())
	}
}
