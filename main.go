package main

import (
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	// "strconv"
	// "strings"
	"./api"
)

func main() {
	http.HandleFunc("/", HttpWrapHandler(get))
	http.HandleFunc("/add", HttpWrapHandler(add))
	/*
		http.HandleFunc("/del", del)
		http.HandleFunc("/edit", edit)
		http.HandleFunc("/login", login)
		http.HandleFunc("/out", out)
	*/
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

type HttpErrorHandler func(w http.ResponseWriter, r *http.Request) error
type HttpHandler func(w http.ResponseWriter, r *http.Request)

func HttpWrapHandler(inHandler HttpErrorHandler) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		err := inHandler(w, r)
		if err != nil {
			if err == api.NoLoginError {
				http.Redirect(w, r, "/login", 302)
			} else {
				fmt.Fprint(w, "<div style=\"color:red\">")
				fmt.Fprint(w, err.Error())
				fmt.Fprint(w, "</div>")
			}
		}
	}
}

func get(w http.ResponseWriter, r *http.Request) error {
	err := api.CheckLogin(r)
	if err != nil {
		return err
	}

	v, err := api.Get("SELECT * FROM book")
	if err != nil {
		return err
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	err = t.Execute(buffer, struct {
		List    []api.Book
		IsLogin bool
	}{
		v,
		true,
	})
	if err != nil {
		return err
	}

	w.Write(buffer.Bytes())
	return nil
}

func add(w http.ResponseWriter, r *http.Request) error {
	err := api.CheckLogin(r)
	if err != nil {
		return err
	}

	if r.Method == "GET" {
		t, err := template.ParseFiles("add.html")
		if err != nil {
			return err
		}

		buffer := bytes.NewBuffer(nil)
		err = t.Execute(w, nil)
		if err != nil {
			return err
		}
		w.Write(buffer.Bytes())
	} else {
		data, err := api.CheckInput(r, map[string]string{
			"username": "string",
			"bname":    "string",
		})
		if err != nil {
			return err
		}

		_, err = api.Add(
			"INSERT book SET Username=?,Bname=?",
			data["username"].(string),
			data["bname"].(string),
		)

		http.Redirect(w, r, "/", 302)
	}
	return nil
}

/*

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
*/
