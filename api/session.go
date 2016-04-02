package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	// "html/template"
)

type SessionStore struct {
	token string
	value map[string]string
	w     http.ResponseWriter
	r     *http.Request
}

func SessionInit(w http.ResponseWriter, r *http.Request) *SessionStore {

	c1, err := r.Cookie("token")

	var s SessionStore

	if err != nil {

		theRandValue := RandString(16)

		_ = Add(
			"INSERT session SET token=?,value=?",
			theRandValue,
			"",
		)

		s = SessionStore{
			token: theRandValue,
			value: map[string]string{},
			w:     w,
			r:     r,
		}
		fmt.Println("无token")
		fmt.Println("theRandValue", theRandValue)
	} else {
		s.getSession("SELECT * FROM session where token=?", c1.Value)

		fmt.Println("有token")

	}

	c := &http.Cookie{
		Name:   "token",
		Value:  s.token,
		Path:   "/",
		MaxAge: 0,
	}
	http.SetCookie(w, c)

	return &s

}

func (st *SessionStore) SessionGet(name string) string {

	single, ok := st.value[name]
	if ok {
		return single
	} else {
		return ""
	}

}

func (st *SessionStore) SessionSet(name string, value string) {

	if st.token == "" {
		c := &http.Cookie{
			Name:   "token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(st.w, c)
		panic(errors.New("非法操作"))
	}

	st.value[name] = value

	fmt.Println("输入了"+name, st.value[name])

}

func (st *SessionStore) SessionDel(name string) {

	_, ok := st.value[name] // 假如key存在,则name = 李四 ，ok = true,否则，ok = false
	if ok {
		delete(st.value, name)
	}

}

func (st *SessionStore) SessionClose() {
	ddb := Db()

	defer ddb.Close()

	stmt, err := ddb.Prepare("update session set value=? where token=?")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(map2json(st.value), st.token)
	if err != nil {
		panic(err)
	}

}

func (st *SessionStore) getSession(query string, args ...interface{}) {

	ddb := Db()
	defer ddb.Close()

	rows, err := ddb.Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		var value string
		err = rows.Scan(&token, &value)
		if err != nil {
			panic(err)
		}

		mapData := json2map(value)

		st.token = token
		st.value = mapData

	}

}

func json2map(jsonString string) map[string]string {

	var content map[string]string
	err := json.Unmarshal([]byte(jsonString), &content)
	if err != nil {
		panic(errors.New("解析JSON失败"))
	}

	return content
}

func map2json(mapData map[string]string) string {

	jsonData, err := json.Marshal(mapData)
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}
