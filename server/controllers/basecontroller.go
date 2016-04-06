package controllers

import (
	// "fmt"
	"bytes"
	"github.com/astaxie/beego"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/web"
	"html/template"
)

type BaseController struct {
	BeegoValidateController
}

func InitRoute(namespace string, target beego.ControllerInterface) {
	InitBeegoVaildateControllerRoute(namespace, target)
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
	this.Ctx.WriteString(string(resultString))
}

func (this *BaseController) AutoRender(returnValue interface{}, viewname string) {
	// result := baseControllerResult{}
	resultError, ok := returnValue.(Exception)
	if ok {
		//带错误码的error
		// result.Code = resultError.GetCode()
		// result.Msg = resultError.GetMessage()
		// result.Data = nil
		this.Ctx.WriteString(resultError.GetMessage())
	} else {
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
		this.Ctx.ResponseWriter.Write(buffer.Bytes())
	}

	// if viewname == "json" {
	// 	this.jsonRender(result)
	// } else if viewname == "html" {
	// 	// this.Ctx.WriteString("OK")
	// 	this.TplName = "/static/index.html"
	// } else {
	// 	panic("不合法的viewName " + viewname)
	// }
}
