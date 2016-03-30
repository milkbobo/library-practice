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

func Db() (*sql.DB, error) {
	return sql.Open("mysql", "root:123@/library?charset=utf8")
}

/*
func Del(id int) {
	ddb := Db()
	defer ddb.Close()

	stmt, err := ddb.Prepare("delete from book where Uid=?")
	CheckErr(err)

	res, err := stmt.Exec(id)
	CheckErr(err)

	_, err = res.RowsAffected()
	CheckErr(err)

}
*/

func Get(query string, args ...interface{}) ([]Book, error) {

	ddb, err := Db()
	if err != nil {
		return nil, err
	}

	defer ddb.Close()

	rows, err := ddb.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	v := []Book{}
	for rows.Next() {
		var uid int
		var username string
		var bname string
		err = rows.Scan(&uid, &username, &bname)
		if err != nil {
			return nil, err
		}

		v = append(v, Book{
			Uid:      uid,
			Username: username,
			Bname:    bname,
		})
	}

	return v, nil
}

/*
func Add(query string, args ...interface{}) int64 {

	ddb := Db()

	defer ddb.Close()

	stmt, err := ddb.Prepare(query)
	CheckErr(err)

	res, err := stmt.Exec(args...)
	CheckErr(err)

	id, err := res.LastInsertId()
	CheckErr(err)

	return id

}

func Edit(query string, args ...interface{}) {

	ddb := Db()
	defer ddb.Close()

	stmt, err := ddb.Prepare(query)
	CheckErr(err)

	res, err := stmt.Exec(args...)
	CheckErr(err)

	_, err = res.RowsAffected()
	CheckErr(err)
}
*/
