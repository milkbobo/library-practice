package main

import (
	"./api"
	// "crypto/rand"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"

	// "encoding/hex"
	"net/http"
	// "strconv"
)

func main() {
	fmt.Println("come on")

	http.HandleFunc("/", HttpWrapHandler(get))

	http.HandleFunc("/add", HttpWrapHandler(add))

	http.HandleFunc("/del", HttpWrapHandler(del))
	http.HandleFunc("/edit", HttpWrapHandler(edit))

	http.HandleFunc("/login", HttpWrapHandler(login))
	http.HandleFunc("/out", HttpWrapHandler(out))

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

type HttpErrorHandler func(w http.ResponseWriter, r *http.Request) error
type HttpHandler func(w http.ResponseWriter, r *http.Request)

func HttpWrapHandler(inHandler HttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				if err == api.NoLoginError {
					http.Redirect(w, r, "/login", 302)
				} else {
					fmt.Println(err)
					fmt.Fprint(w, "<div style=\"color:red\">")
					fmt.Fprintf(w, "%v", err)
					fmt.Fprint(w, "</div>")
				}
			}
			// panic(err)
		}()
		inHandler(w, r)

		fmt.Println("come out")
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	api.CheckLogin(w, r)

	v := api.Get("SELECT * FROM book")

	result := api.TemplateOutput("index.html", struct {
		List    []api.Book
		IsLogin bool
	}{
		v,
		true,
	})

	w.Write(result)
}

func add(w http.ResponseWriter, r *http.Request) {
	api.CheckLogin(w, r)

	if r.Method == "GET" {
		result := api.TemplateOutput("add.html", nil)

		w.Write(result)
	} else {
		data := api.CheckInput(r, map[string]string{
			"username": "string",
			"bname":    "string",
		})

		_ = api.Add(
			"INSERT book SET Username=?,Bname=?",
			data["username"].(string),
			data["bname"].(string),
		)

		http.Redirect(w, r, "/", 302)
	}

}

func del(w http.ResponseWriter, r *http.Request) {
	api.CheckLogin(w, r)

	data := api.CheckInput(r, map[string]string{
		"id": "int",
	})

	v := api.Get("SELECT * FROM book where Uid=?", data["id"].(int))

	if len(v) == 0 {
		panic(errors.New("不存在该数据"))
	}

	api.Del(data["id"].(int))

	http.Redirect(w, r, "/", 302)

}

func edit(w http.ResponseWriter, r *http.Request) {
	api.CheckLogin(w, r)

	if r.Method == "GET" {
		data := api.CheckInput(r, map[string]string{
			"id": "int",
		})

		v := api.Get("SELECT * FROM book where Uid=?", data["id"].(int))

		if len(v) == 0 {
			panic(errors.New("不存在该数据"))
		}

		result := api.TemplateOutput("edit.html", v[0])

		w.Write(result)
	} else {
		data := api.CheckInput(r, map[string]string{
			"id":    "int",
			"bname": "string",
		})

		api.Edit(
			"update book set Bname=? where uid=?",
			data["bname"].(string),
			data["id"].(int),
		)

		http.Redirect(w, r, "/", 302)
	}

}

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		result := api.TemplateOutput("login.html", nil)

		w.Write(result)
	} else {
		data := api.CheckInput(r, map[string]string{
			"username": "string",
			"password": "string",
		})

		if data["username"] != "admin" {
			panic(errors.New("账号错误"))
		}
		if data["password"] != "admin" {
			panic(errors.New("密码错误"))
		}

		s := api.SessionStore{}

		s.SessionSet(w, r, "admin")
		http.Redirect(w, r, "/", 302)
	}

}

func out(w http.ResponseWriter, r *http.Request) {

	s := api.SessionStore{}
	s.SessionDestroy(w, r)

	http.Redirect(w, r, "/", 302)

}
