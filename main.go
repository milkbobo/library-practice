package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	// "strconv"
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

	v := api.Get("SELECT * FROM book")

	fmt.Printf("%v\n", v)

	IsLogin := false

	c1, err := r.Cookie("username")

	fmt.Println("c1", c1)

	if c1 != nil {
		IsLogin = true
	}

	t, _ := template.ParseFiles("index.html")
	err = t.Execute(w, struct {
		List    []api.Book
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

		username := template.HTMLEscapeString(r.Form.Get("username"))
		bname := template.HTMLEscapeString(r.Form.Get("bname"))
		fmt.Println(username, bname)

		id := api.Add("INSERT book SET Username=?,Bname=?", username, bname)

		fmt.Println(id)

		http.Redirect(w, r, "/", 302)
		return

	}

}

func del(w http.ResponseWriter, r *http.Request) {

	id := api.CheckId(w, r)

	if id == 0 {
		return
	}

	api.Del(id)

	http.Redirect(w, r, "/", 302)
	return

}

func edit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	idd := template.HTMLEscapeString(r.Form.Get("id"))
	// username := template.HTMLEscapeString(r.Form.Get("username"))
	bname := template.HTMLEscapeString(r.Form.Get("bname"))
	fmt.Printf("the.idd->[%v]\n", idd)
	fmt.Printf("the.bname->[%v]\n", bname)

	c1 := api.CheckLogin(w, r)

	ids := api.CheckId(w, r)

	// rows, err := ddb.Query("SELECT * FROM book where Uid=?", ids)

	v := api.Get("SELECT * FROM book where Uid=?", ids)

	if len(v) == 0 {
		fmt.Fprint(w, "非法操作，请返回重试")
		return
	}

	fmt.Printf("%v\n", v)

	if r.Method == "GET" {

		t, _ := template.ParseFiles("edit.html")
		err := t.Execute(w, v[0])
		api.CheckErr(err)
	} else {

		api.Edit("update book set Username=?,Bname=? where uid=?", c1.Value, bname, idd)

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
