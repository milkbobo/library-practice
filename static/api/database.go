package api

import (
	"database/sql"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Bid      int
	Username string
	Bname    string
}
type Userinfo struct {
	Uid      int
	Username string
	Password string
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
		var bid int
		var username string
		var bname string
		err = rows.Scan(&bid, &username, &bname)
		if err != nil {
			panic(err)
		}

		v = append(v, Book{
			Bid:      bid,
			Username: username,
			Bname:    bname,
		})
	}

	return v
}

func Del(id int) {
	ddb := Db()

	defer ddb.Close()

	stmt, err := ddb.Prepare("delete from book where Bid=?")
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

func GetUserinfo(query string, args ...interface{}) []Userinfo {

	ddb := Db()
	defer ddb.Close()

	rows, err := ddb.Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	v := []Userinfo{}
	for rows.Next() {
		var uid int
		var username string
		var password string
		err = rows.Scan(&uid, &username, &password)
		if err != nil {
			panic(err)
		}

		v = append(v, Userinfo{
			Uid:      uid,
			Username: username,
			Password: password,
		})
	}

	return v
}
