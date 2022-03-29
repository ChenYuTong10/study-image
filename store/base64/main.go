package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

var conn *sql.DB

func init() {
	// connect the database
	var err error
	conn, err = sql.Open("mysql", os.Getenv("MYSQL_DEV_URL"))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	// set cors header
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get image information
	name := r.PostFormValue("name")
	code := r.PostFormValue("code")

	// you can have a check here

	// insert into the database
	sql := `
		INSERT INTO t_base64_image
		(name, code)
		VALUES
		(?, ?);
	`
	_, err := conn.Exec(sql, name, code)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "db fail", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":9090", nil)
}