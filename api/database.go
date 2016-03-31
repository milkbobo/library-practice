package api

import (
	"database/sql"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Uid      int
	Username string
	Bname    string
}

type Session struct {
	token string
	value string
}

func Db() *sql.DB {
	db, err := sql.Open("mysql", "root:1@/library?charset=utf8")
	if err != nil {
		panic(err)
	}
	return db
}

func Get(query string, args ...interface{}) []Book {

	ddb := Db()
	defer ddb.Close()

	rows, err := ddb.Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	v := []Book{}
	for rows.Next() {
		var uid int
		var username string
		var bname string
		err = rows.Scan(&uid, &username, &bname)
		if err != nil {
			panic(err)
		}

		v = append(v, Book{
			Uid:      uid,
			Username: username,
			Bname:    bname,
		})
	}

	return v
}

func GetSession(query string, args ...interface{}) []Session {

	ddb := Db()
	defer ddb.Close()

	rows, err := ddb.Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	v := []Session{}
	for rows.Next() {
		var token string
		var value string
		err = rows.Scan(&token, &value)
		if err != nil {
			panic(err)
		}

		v = append(v, Session{
			token: token,
			value: value,
		})
	}

	return v
}

func Del(id int) {
	ddb := Db()

	defer ddb.Close()

	stmt, err := ddb.Prepare("delete from book where Uid=?")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		panic(err)
	}

}

func Add(query string, args ...interface{}) int64 {

	ddb := Db()

	defer ddb.Close()

	stmt, err := ddb.Prepare(query)
	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return id

}

func Edit(query string, args ...interface{}) {

	ddb := Db()

	defer ddb.Close()

	stmt, err := ddb.Prepare(query)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		panic(err)
	}

}
