package api

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strconv"
)

var NoLoginError = errors.New("账号未登录")

func CheckLogin(r *http.Request) error {
	c1, err := r.Cookie("username")
	if err != nil || c1.Value != "admin" {
		return NoLoginError
	}
	return nil
}

func CheckInput(r *http.Request, inputFilter map[string]string) (map[string]interface{}, error) {
	r.ParseForm()
	result := map[string]interface{}{}
	for key, format := range inputFilter {
		singleData := r.Form.Get(key)
		if singleData == "" {
			return nil, errors.New("缺少参数" + key)
		}
		var singleResult interface{}
		if format == "string" {
			singleResult = singleData
		} else if format == "int" {
			var err error
			singleResult, err = strconv.Atoi(singleData)
			if err != nil {
				return nil, errors.New(key + "参数不是合法的整数:[" + singleData + "]")
			}
		} else {
			return nil, errors.New("不合法的format" + format)
		}
		result[key] = singleResult
	}
	return result, nil
}

func TemplateOutput(filename string, data interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	t, err := template.ParseFiles(filename)
	if err != nil {
		return nil, err
	}

	err = t.Execute(buffer, data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
