package main

import (
	"./api"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
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

	result, err := api.TemplateOutput("index.html", struct {
		List    []api.Book
		IsLogin bool
	}{
		v,
		true,
	})
	if err != nil {
		return err
	}

	w.Write(result)
	return nil
}

func add(w http.ResponseWriter, r *http.Request) error {
	err := api.CheckLogin(r)
	if err != nil {
		return err
	}

	if r.Method == "GET" {
		result, err := api.TemplateOutput("add.html", nil)
		if err != nil {
			return err
		}

		w.Write(result)
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

func del(w http.ResponseWriter, r *http.Request) error {
	err := api.CheckLogin(r)
	if err != nil {
		return err
	}

	data, err := api.CheckInput(r, map[string]string{
		"id": "int",
	})
	if err != nil {
		return err
	}

	v, err := api.Get("SELECT * FROM book where Uid=?", data["id"].(int))
	if err != nil {
		return err
	}
	if len(v) == 0 {
		return errors.New("不存在该数据")
	}

	err = api.Del(data["id"].(int))
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", 302)
	return nil
}

func edit(w http.ResponseWriter, r *http.Request) error {
	err := api.CheckLogin(r)
	if err != nil {
		return err
	}

	if r.Method == "GET" {
		data, err := api.CheckInput(r, map[string]string{
			"id": "int",
		})
		if err != nil {
			return err
		}

		v, err := api.Get("SELECT * FROM book where Uid=?", data["id"].(int))
		if err != nil {
			return err
		}
		if len(v) == 0 {
			return errors.New("不存在该数据")
		}

		result, err := api.TemplateOutput("edit.html", v[0])
		if err != nil {
			return err
		}

		w.Write(result)
	} else {
		data, err := api.CheckInput(r, map[string]string{
			"id":    "int",
			"bname": "string",
		})
		if err != nil {
			return err
		}

		err = api.Edit(
			"update book set Bname=? where uid=?",
			data["bname"].(string),
			data["id"].(int),
		)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/", 302)
	}
	return nil

}

func login(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		result, err := api.TemplateOutput("login.html", nil)
		if err != nil {
			return err
		}

		w.Write(result)
	} else {
		data, err := api.CheckInput(r, map[string]string{
			"username": "string",
			"password": "string",
		})
		if err != nil {
			return err
		}

		if data["username"] != "admin" {
			return errors.New("账号错误")
		}
		if data["password"] != "admin" {
			return errors.New("密码错误")
		}

		c := &http.Cookie{
			Name:   "username",
			Value:  "admin",
			Path:   "/",
			MaxAge: 120,
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, "/", 302)
	}
	return nil
}

func out(w http.ResponseWriter, r *http.Request) error {
	c := &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	http.Redirect(w, r, "/", 302)
	return nil
}
