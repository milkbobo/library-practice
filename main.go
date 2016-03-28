package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    // "strings"
     _ "github.com/go-sql-driver/mysql"
    "database/sql"
)




func main() {    
    http.HandleFunc("/", get)     
    http.HandleFunc("/add",add)
    http.HandleFunc("/del",del)
    http.HandleFunc("/edit",edit)
    err := http.ListenAndServe(":9090", nil) 
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

}

func db() (*sql.DB) {
	dbs, err := sql.Open("mysql", "root:milkbobo@/library?charset=utf8")
	checkErr(err)
	return dbs
}

func get(w http.ResponseWriter, r *http.Request ){

	ddb := db()

	defer ddb.Close()

	rows, err := ddb.Query("SELECT * FROM book")
    checkErr(err)

    type dd struct{
    	Uid int
    	Username string
    	Bname string
    }

    v := []dd{}

    for rows.Next() {
        var uid int
        var username string
        var bname string
        err = rows.Scan(&uid, &username, &bname,)
        checkErr(err)
        fmt.Println(uid)
        fmt.Println(username)
        fmt.Println(bname)

        v = append(v,dd{
        	Uid:uid,
        	Username:username,
        	Bname:bname,
        })
    }

    	fmt.Printf("%v\n", v)

        t, _ := template.ParseFiles("index.html")
        err=t.Execute(w,v)
        if err != nil{

        fmt.Println(err.Error())
        }

}

func add(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fmt.Println("method:", r.Method) 
    if r.Method == "GET" {
        t, _ := template.ParseFiles("add.html")
        t.Execute(w, nil)
    } else {

		username :=r.Form.Get("username")
		bname := r.Form.Get("bname")
		fmt.Println(username,bname)

	    ddb := db()

		defer ddb.Close()

		stmt, err := ddb.Prepare("INSERT book SET Username=?,Bname=?")
	    checkErr(err)

	    res, err := stmt.Exec(username, bname)
	    checkErr(err)

	    id, err := res.LastInsertId()
	    checkErr(err)

	    fmt.Println(id)

	    http.Redirect(w,r,"/",302)
	    return


    }



}

func del(w http.ResponseWriter, r *http.Request){
	ddb := db()

	defer ddb.Close()

	r.ParseForm()  

	if len(r.Form["id"]) > 0 {

	    //删除数据
    stmt, err := ddb.Prepare("delete from book where Uid=?")
    checkErr(err)

    res, err := stmt.Exec(r.Form["id"][0])
    checkErr(err)

    _, err = res.RowsAffected()
    checkErr(err)

    http.Redirect(w,r,"/",302)

	   
	}  
}

func edit(w http.ResponseWriter, r *http.Request){

	ddb := db()

	defer ddb.Close()

	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法

    if r.Method == "GET" {

    	rows, err := ddb.Query("SELECT * FROM book where Uid="+r.Form["id"][0])
  	 	checkErr(err)




		type dd struct{
			Uid int
			Username string
			Bname string
		}

		v :=dd{}

		for rows.Next() {
		    var uid int
		    var username string
		    var bname string
		    err = rows.Scan(&uid, &username, &bname,)
		    checkErr(err)
		    fmt.Println(uid)
		    fmt.Println(username)
		    fmt.Println(bname)

		    v = dd{
		    	Uid:uid,
		    	Username:username,
		    	Bname:bname,
		    }
		}


   		fmt.Printf("%v\n", v)

        t, _ := template.ParseFiles("edit.html")
        err = t.Execute(w, v)
        checkErr(err)
    } else {
		id :=r.Form.Get("id")
		username :=r.Form.Get("username")
		bname := r.Form.Get("bname")
		fmt.Println(id,username,bname)


	    stmt, err := ddb.Prepare("update book set Username=?,Bname=? where uid=?")
	    checkErr(err)

	    res, err := stmt.Exec(username,bname, id)
	    checkErr(err)

	    _, err = res.RowsAffected()
	    checkErr(err)

	    fmt.Println(id)

	    http.Redirect(w,r,"/",302)
	    return


    }



}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}