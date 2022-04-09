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
	"strconv"
	"time"
)

var (
	conn     *sql.DB
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

	// get image binary content and other parameters
	userId := r.PostFormValue("userId")
	index := r.PostFormValue("index")
	chunk, _, _ := r.FormFile("chunk")

	// make a directory to store chunk from the same image
	dirPath := fmt.Sprintf("%v\\%v", rootPath, userId)

	if err := os.MkdirAll(dirPath, 0666); err != nil {
		log.Println(err.Error())
		http.Error(w, "make dir fail", http.StatusInternalServerError)
		return
	}

	// create a chunk
	chunkPath := fmt.Sprintf("%v\\%v", dirPath, index)

	local, err := os.Create(chunkPath)
	defer local.Close()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "create chunk fail", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(local, chunk)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "read chunk fail", http.StatusInternalServerError)
		return
	}

	// get the delay of the interface
	delay := fmt.Sprintf("%s", time.Since(startTime))

	// json serialize
	strResp, _ := json.Marshal(&resp{StatusCode: 2000, Message: "success", Delay: delay})

	w.Write(strResp)
}

func merge(w http.ResponseWriter, r *http.Request) {

	// get merge information
	userId := r.PostFormValue("userId")
	imageExt := r.PostFormValue("imageExt")
	strChunkNum := r.PostFormValue("chunkNum")

	// transfer `string` type to `int64`
	chunkNum, _ := strconv.ParseInt(strChunkNum, 10, 64)

	// get an unique id
	id, _ := getUniqueId()

	imageName := fmt.Sprintf("%v%v", id, imageExt)
	absPath := fmt.Sprintf("%v\\%v", rootPath, imageName)

	// create a new file to store the merged image
	local, err := os.OpenFile(absPath,  os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer local.Close()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "create file fail", http.StatusInternalServerError)
		return
	}

	// merge all chunks in the directory named `userId`
	chunkDirPath := fmt.Sprintf("%v\\%v", rootPath, userId)
	for i := int64(0); i < chunkNum; i++ {

		chunkPath := fmt.Sprintf("%v\\%v", chunkDirPath, i)

		chunk, err := os.Open(chunkPath)
		if err != nil {
			chunk.Close()
			log.Println(err.Error())
			http.Error(w, "open chunk fail", http.StatusInternalServerError)
			return
		}

		// loop to read the chunk
		for {
			buffer := make([]byte, 1024)

			_, err = chunk.Read(buffer)
			if err == io.EOF {
				break
			}

			_, err = local.Write(buffer)
			if err != nil {
				chunk.Close()
				log.Println(err.Error())
				http.Error(w, "write fail", http.StatusInternalServerError)
				return
			}
		}

		chunk.Close()
	}

	// insert image info into the database
	insertSQL := `
		INSERT INTO slice_upload_image
		(path)
		VALUES
		(?);
	`
	_, err = conn.Exec(insertSQL, imageName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "db fail", http.StatusInternalServerError)
		return
	}

	// remove all the chunks
	defer func() {
		err = os.RemoveAll(chunkDirPath)
		if err != nil {
			log.Println(err.Error())
		}
	}()

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
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/merge", merge)
	http.Handle("/", staticFileServer())
	http.ListenAndServe(":9090", nil)
}
