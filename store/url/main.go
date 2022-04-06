package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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
	remote, header, _ := r.FormFile("image")

	// you can have a check here
	if header.Size > 10*2<<20 {
		http.Error(w, "too big", http.StatusBadRequest)
		return
	}

	// ensure you have the unique image name to avoid the name conflict
	id, _ := getUniqueId()
	extension := path.Ext(header.Filename)

	imageName := fmt.Sprintf("%v%v", id, extension)
	absPath := fmt.Sprintf("%v/%v", rootPath, imageName)

	// store the image to local
	local, err := os.Create(absPath)
	defer local.Close() // it must be closed when the image finish reading
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "create fail", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(local, remote)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "read fail", http.StatusInternalServerError)
		return
	}

	// insert into the database
	sql := `
		INSERT INTO url_store_image
		(path)
		VALUES
		(?);
	`
	_, err = conn.Exec(sql, absPath)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "db fail", http.StatusInternalServerError)
		return
	}

	// get the delay of the interface
	delay := fmt.Sprintf("%s", time.Since(startTime))

	// json serialize
	strResp, _ := json.Marshal(&resp{StatusCode: 2000, Message: imageName, Delay: delay})

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