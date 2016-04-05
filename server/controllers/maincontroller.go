package controllers

import (
	. "github.com/fishedee/language"
)

type MainController struct {
	BaseController
}

func (this *MainController) Test_Json() interface{} {
	Throw(1, "你输入缺少了一个参数")
	return "12xxx"
}
