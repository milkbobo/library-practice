package api

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var NoLoginError = errors.New("账号未登录")

func CheckLogin(w http.ResponseWriter, r *http.Request) {

	s := SessionInit(w, r)
	defer s.SessionClose()
	v := s.SessionGet("username")

	fmt.Println(v)

	if v == "" {
		panic(NoLoginError)
	}

	// if v != "admin" {
	// 	panic(NoLoginError)
	// }
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

func RandString(num int) string {
	//生成随机数
	k := make([]byte, num)
	if _, err := rand.Read(k); err != nil {
		panic(err)
	}

	return hex.EncodeToString(k)
}

func CheckCsrf(w http.ResponseWriter, r *http.Request) {
	data := CheckInput(r, map[string]string{
		"csrf": "string",
	})

	c1, err := r.Cookie("token")

	if err != nil {
		c := &http.Cookie{
			Name:   "token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, c)
		panic(errors.New("无效登陆，请重新登陆"))
	}

	fmt.Println("csrf", data["csrf"])
	fmt.Println("c1.Value", c1.Value)

	if data["csrf"] != c1.Value+"edward" {
		panic(errors.New("我被csrf攻击啦！"))
	}
}
