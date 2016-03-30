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
	return sql.Open("mysql", "root:1@/library?charset=utf8")
}

func Del(id int) error {
	ddb, err := Db()
	if err != nil {
		return err
	}
	defer ddb.Close()

	stmt, err := ddb.Prepare("delete from book where Uid=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

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

func Add(query string, args ...interface{}) (int64, error) {

	ddb, err := Db()
	if err != nil {
		return 0, err
	}

	defer ddb.Close()

	stmt, err := ddb.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil

}

func Edit(query string, args ...interface{}) error {

	ddb, err := Db()
	if err != nil {
		return err
	}
	defer ddb.Close()

	stmt, err := ddb.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}
