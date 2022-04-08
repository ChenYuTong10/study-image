package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	conn 	*sql.DB
	rootPath string
)

type resp struct {
	StatusCode int64  `json:"statusCode"`
	Message    string `json:"message"`
	Delay      string `json:"delay"`
}

func init() {
	// connect the database
	var err error
	conn, err = sql.Open("mysql", os.Getenv("MYSQL_DEV_URL"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	// get the root path to store the image file
	rootPath, err = os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	// time tick
	startTime := time.Now()

	// set cors header
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get image information
	name := r.PostFormValue("name")
	code := r.PostFormValue("code")

	// you can have a check here

	// insert into the database
	sql := `
		INSERT INTO base64_store_image
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

	// get the delay of the interface
	delay := fmt.Sprintf("%s", time.Since(startTime))

	// json serialize
	strResp, _ := json.Marshal(&resp{StatusCode: 2000, Delay: delay})

	w.Write(strResp)
}

func staticFileServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set cors header
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// you can also have a check here

		// appoint a directory to store the resources
		http.FileServer(http.Dir(rootPath)).ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/upload", upload)
	http.Handle("/", staticFileServer())
	http.ListenAndServe(":9090", nil)
}