package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	// "html/template"
)

type SessionStore struct {
	token string
	value string
	w     http.ResponseWriter
	r     *http.Request
}

func SessionInit(w http.ResponseWriter, r *http.Request) *SessionStore {

	c1, err := r.Cookie("token")

	var s SessionStore

	if err != nil {
		k := make([]byte, 16)
		if _, err := rand.Read(k); err != nil {
			panic(err)
		}

		theRandValue := hex.EncodeToString(k)

		s = SessionStore{
			token: theRandValue,
			value: "",
			w:     w,
			r:     r,
		}
		fmt.Println("无token")
	} else {
		v := GetSession("SELECT * FROM session where token=?", c1.Value)

		s = SessionStore{
			token: v[0].token,
			value: v[0].value,
			w:     w,
			r:     r,
		}

		fmt.Println("有token")

	}

	return &s

}

func (st *SessionStore) SessionSet(name string, value string) {

	st.Input(name, value)

	_ = Add(
		"INSERT session SET token=?,value=?",
		st.token,
		st.value,
	)

	c := &http.Cookie{
		Name:   "token",
		Value:  st.token,
		Path:   "/",
		MaxAge: 0,
	}
	http.SetCookie(st.w, c)
}

func (st *SessionStore) SessionGet(name string) string {

	if st.value == "" {

		fmt.Println("为空")
		return ""

	} else {
		fmt.Println("有东西")
		fmt.Println(st.value)
		content := make(map[string]string)
		err := json.Unmarshal([]byte(st.value), &content)
		if err != nil {
			panic(errors.New("解析JSON失败"))
		}
		fmt.Println(content)
		singleValue := content[name]
		return singleValue
	}

}

func (st *SessionStore) SessionDestroy(name string) {
	if st.value != "" {
		content := make(map[string]string)
		err := json.Unmarshal([]byte(st.value), &content)
		if err != nil {
			panic(err)
		}

		_, ok := content[name] // 假如key存在,则name = 李四 ，ok = true,否则，ok = false
		if ok {
			delete(content, name)
		}
	}
}

func (st *SessionStore) SessionClose() {

	Close(st.token)

	c := &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	fmt.Println(st.w, c)

	http.SetCookie(st.w, c)

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

func Close(token string) {
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

func (st *SessionStore) Input(name string, value string) {

	content := make(map[string]string)

	if st.value == "" {
		content[name] = value
	} else {
		err := json.Unmarshal([]byte(st.value), &content)
		if err != nil {
			panic(err)
		}
		content[name] = value
	}

	jsonVelue, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	st.value = string(jsonVelue)

}
