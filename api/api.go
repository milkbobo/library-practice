package api

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var NoLoginError = errors.New("账号未登录")

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	c1, err := r.Cookie("token")

	if err != nil {
		http.Redirect(w, r, "/login", 302)
		panic(err)
	}

	token := c1.Value
	fmt.Println("token", token)

	v := GetSession("SELECT * FROM session where token=?", token)

	fmt.Println("v", v)

	if len(v) == 0 {
		panic(errors.New("非法登陆"))
	}

	if err != nil || v[0].value != "admin" {
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
