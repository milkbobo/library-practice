package api

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strconv"
)

var NoLoginError = errors.New("账号未登录")

func CheckLogin(r *http.Request) {
	c1, err := r.Cookie("username")
	if err != nil || c1.Value != "admin" {
		panic(NoLoginError)
	}
}

func CheckInput(r *http.Request, inputFilter map[string]string) map[string]interface{} {
	r.ParseForm()
	result := map[string]interface{}{}
	for key, format := range inputFilter {
		singleData := r.Form.Get(key)
		if singleData == "" {
			panic(errors.New("缺少参数" + key))
		}
		var singleResult interface{}
		if format == "string" {
			singleResult = singleData
		} else if format == "int" {
			var err error
			singleResult, err = strconv.Atoi(singleData)
			if err != nil {
				panic(errors.New(key + "参数不是合法的整数:[" + singleData + "]"))
			}
		} else {
			panic(errors.New("不合法的format" + format))
		}
		result[key] = singleResult
	}
	return result
}

func TemplateOutput(filename string, data interface{}) []byte {
	buffer := bytes.NewBuffer(nil)
	t, err := template.ParseFiles(filename)
	if err != nil {
		panic(err)
	}

	err = t.Execute(buffer, data)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}
