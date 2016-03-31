package api

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	// "html/template"
)

type SessionStore struct {
	token string
	value string
}

func (st *SessionStore) SessionInit() {
	k := make([]byte, 16)
	if _, err := rand.Read(k); err != nil {
		panic(err)
	}

	theRandValue := hex.EncodeToString(k)
	st.token = theRandValue
}

func (st *SessionStore) SessionSet(w http.ResponseWriter, r *http.Request, value string) {

	c1, err := r.Cookie("token")

	if err != nil {
		st.SessionInit()
	} else {
		st.token = c1.Value
	}

	_ = Add(
		"INSERT session SET token=?,value=?",
		st.token,
		value,
	)

	c := &http.Cookie{
		Name:   "token",
		Value:  st.token,
		Path:   "/",
		MaxAge: 0,
	}
	http.SetCookie(w, c)
}

func (st *SessionStore) SessionGet(w http.ResponseWriter, r *http.Request) []SessionStore {
	c1, err := r.Cookie("token")

	if err != nil {
		panic(NoLoginError)
	}

	token := c1.Value
	fmt.Println("token", token)

	v := GetSession("SELECT * FROM session where token=?", token)

	fmt.Println("v", v)

	if len(v) == 0 {
		panic(errors.New("非法登陆"))
	}

	return v
}

func (st *SessionStore) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	c1, err := r.Cookie("token")

	if err == nil {
		token := c1.Value
		DelSession(token)
	}

}

func GetSession(query string, args ...interface{}) []SessionStore {

	ddb := Db()
	defer ddb.Close()

	rows, err := ddb.Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	v := []SessionStore{}
	for rows.Next() {
		var token string
		var value string
		err = rows.Scan(&token, &value)
		if err != nil {
			panic(err)
		}

		v = append(v, SessionStore{
			token: token,
			value: value,
		})
	}

	return v
}

func DelSession(token string) {
	ddb := Db()

	defer ddb.Close()

	stmt, err := ddb.Prepare("delete from session where token=?")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(token)
	if err != nil {
		panic(err)
	}

}
