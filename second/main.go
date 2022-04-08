package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	conn 	*sql.DB
	rootPath string
)

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

func check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get the md5 value of the image
	md5 := r.PostFormValue("md5")

	// generally, you should check the validity of parameters

	// find the same value in database
	readSQL := `
		SELECT path
		FROM second_transmit_image
		WHERE md5 = ?;
	`
	results := conn.QueryRow(readSQL, md5)

	var imageName string
	err := results.Scan(&imageName)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		http.Error(w, "db fail", http.StatusInternalServerError)
		return
	}

	// exist the same image
	if len(imageName) > 0 {
		w.Write([]byte(imageName))
		return
	}

	// you can write back in any way instead of the simple way like me
	w.Write([]byte("no"))

}

// upload is roughly same as the `upload` function in `/store/url`
func upload(w http.ResponseWriter, r *http.Request) {
	// set cors header
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get image information
	md5 := r.PostFormValue("md5")
	remote, header, _ := r.FormFile("image")

	// get a unique image id
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
	insertSQL := `
		INSERT INTO second_transmit_image
		(md5, path)
		VALUES
		(?, ?);
	`
	_, err = conn.Exec(insertSQL, md5, imageName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "db fail", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(imageName))
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
	http.HandleFunc("/check", check)
	http.HandleFunc("/upload", upload)
	http.Handle("/", staticFileServer())
	http.ListenAndServe(":9090", nil)
}