package api

import (
	"errors"
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

/*
type dd struct {
	Uid      int
	Username string
	Bname    string
}

func CheckId(w http.ResponseWriter, r *http.Request) int {
	ddb := Db()

	defer ddb.Close()

	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法

	if len(r.Form["id"]) <= 0 {
		fmt.Fprint(w, "请输入id参数")
		return 0
	}

	id := r.Form["id"][0]

	ids, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprint(w, "输入id参数错误，请返回重试！")
		return 0
	}

	fmt.Println("ids", ids)

	v := Get("SELECT * FROM book where Uid=?", ids)

	if len(v) == 0 {
		fmt.Fprint(w, "非法操作，请返回重试")
		return 0
	}

	fmt.Printf("%v\n", v)

	c1 := CheckLogin(w, r)

	fmt.Printf("test,%#v\n", c1)

	if c1 == nil {
		return 0
	}

	if c1.Value != v[0].Username {
		fmt.Fprint(w, "你不是该拥有者，不能删除或修改")
		return 0
	}

	return ids

}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

*/
