package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	// "strings"
	"./api"
)

type dd struct {
	Uid      int
	Username string
	Bname    string
}

func main() {
	http.HandleFunc("/", get)
	http.HandleFunc("/add", add)
	http.HandleFunc("/del", del)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/login", login)
	http.HandleFunc("/out", out)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func get(w http.ResponseWriter, r *http.Request) {

	ddb := api.Db()

	defer ddb.Close()

	rows, err := ddb.Query("SELECT * FROM book")
	api.CheckErr(err)

	v := []dd{}

	for rows.Next() {
		var uid int
		var username string
		var bname string
		err = rows.Scan(&uid, &username, &bname)
		api.CheckErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(bname)

		v = append(v, dd{
			Uid:      uid,
			Username: username,
			Bname:    bname,
		})
	}

	fmt.Printf("%v\n", v)

	IsLogin := false

	c1, err := r.Cookie("username")

	fmt.Println("c1", c1)

	if c1 != nil {
		IsLogin = true
	}

	t, _ := template.ParseFiles("index.html")
	err = t.Execute(w, struct {
		List    []dd
		IsLogin bool
	}{
		v,
		IsLogin,
	})

	if err != nil {

		fmt.Println(err.Error())
	}

}

func add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("add.html")
		t.Execute(w, nil)
	} else {

		username := r.Form.Get("username")
		bname := r.Form.Get("bname")
		fmt.Println(username, bname)

		ddb := api.Db()

		defer ddb.Close()

		stmt, err := ddb.Prepare("INSERT book SET Username=?,Bname=?")
		api.CheckErr(err)

		res, err := stmt.Exec(username, bname)
		api.CheckErr(err)

		id, err := res.LastInsertId()
		api.CheckErr(err)

		fmt.Println(id)

		http.Redirect(w, r, "/", 302)
		return

	}

}

func del(w http.ResponseWriter, r *http.Request) {
	ddb := api.Db()
	defer ddb.Close()

	id := api.CheckId(w, r)

	if id == 0 {
		return
	}

	stmt, err := ddb.Prepare("delete from book where Uid=?")
	api.CheckErr(err)

	res, err := stmt.Exec(id)
	api.CheckErr(err)

	_, err = res.RowsAffected()
	api.CheckErr(err)

	http.Redirect(w, r, "/", 302)

}

func edit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	idd := template.HTMLEscapeString(r.Form.Get("id"))
	// username := template.HTMLEscapeString(r.Form.Get("username"))
	bname := template.HTMLEscapeString(r.Form.Get("bname"))
	fmt.Printf("the.idd->[%v]\n", idd)
	fmt.Printf("the.bname->[%v]\n", bname)

	c1 := api.CheckLogin(w, r)

	ddb := api.Db()

	defer ddb.Close()

	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法

	if len(r.Form["id"]) <= 0 {
		fmt.Fprint(w, "请输入id参数")
		return
	}

	id := r.Form["id"][0]

	ids, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprint(w, "输入id参数错误，请返回重试！")
		return
	}

	fmt.Println("ids", ids)

	rows, err := ddb.Query("SELECT * FROM book where Uid=?", ids)
	api.CheckErr(err)

	defer rows.Close()

	v := []dd{}

	err = rows.Err()
	api.CheckErr(err)

	for rows.Next() {
		var uid int
		var username string
		var bname string
		err = rows.Scan(&uid, &username, &bname)
		api.CheckErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(bname)

		v = append(v, dd{
			Uid:      uid,
			Username: username,
			Bname:    bname,
		})
	}

	if len(v) == 0 {
		fmt.Fprint(w, "非法操作，请返回重试")
		return
	}

	fmt.Printf("%v\n", v)

	if r.Method == "GET" {

		t, _ := template.ParseFiles("edit.html")
		err = t.Execute(w, v[0])
		api.CheckErr(err)
	} else {

		fmt.Println("输入进来的", id, bname)

		stmt, err := ddb.Prepare("update book set Username=?,Bname=? where uid=?")
		api.CheckErr(err)

		res, err := stmt.Exec(c1.Value, bname, idd)
		api.CheckErr(err)

		_, err = res.RowsAffected()
		api.CheckErr(err)

		fmt.Println(id)

		http.Redirect(w, r, "/", 302)
		return

	}

}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if username != "admin" {
			fmt.Fprint(w, "账号或密码错误")
			return
		}

		if password != "admin" {
			fmt.Fprint(w, "账号或密码错误")
			return
		}

		c := &http.Cookie{
			Name:  "username",
			Value: "admin",
			Path:  "/",
			// Domain: "localhost",
			MaxAge: 120,
		}
		http.SetCookie(w, c)

		http.Redirect(w, r, "/", 302)
		return
	}
}

func out(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:  "username",
		Value: "",
		Path:  "/",
		// Domain: "localhost",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	http.Redirect(w, r, "/", 302)
	return
}

func checkLogin(w http.ResponseWriter, r *http.Request) *http.Cookie {

	c1, err := r.Cookie("username")

	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return nil
	}
	return c1

}

func checkId(w http.ResponseWriter, r *http.Request) int {
	ddb := api.Db()

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

	rows, err := ddb.Query("SELECT * FROM book where Uid=?", ids)
	api.CheckErr(err)

	defer rows.Close()

	v := []dd{}

	err = rows.Err()
	api.CheckErr(err)

	for rows.Next() {
		var uid int
		var username string
		var bname string
		err = rows.Scan(&uid, &username, &bname)
		api.CheckErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(bname)

		v = append(v, dd{
			Uid:      uid,
			Username: username,
			Bname:    bname,
		})
	}

	if len(v) == 0 {
		fmt.Fprint(w, "非法操作，请返回重试")
		return 0
	}

	fmt.Printf("%v\n", v)
	c1 := api.CheckLogin(w, r)

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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
