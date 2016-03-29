package api

import (
	"fmt"

	"net/http"
	"strconv"
)

type dd struct {
	Uid      int
	Username string
	Bname    string
}

func CheckLogin(w http.ResponseWriter, r *http.Request) *http.Cookie {

	c1, err := r.Cookie("username")

	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return nil
	}
	return c1

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
