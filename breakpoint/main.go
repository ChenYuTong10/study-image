package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ctx      context.Context
	rdb      *redis.Client
	rootPath string
)

// HistoryKey the prefix key in redis
const HistoryKey = "history:"

func init() {
	var err error

	// get the root path
	rootPath, err = os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	ctx = context.Background()

	// connect to the redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err = rdb.Ping(ctx).Err(); err != nil {
		log.Fatalln(err.Error())
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	// set cors header
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get the file information
	remote, header, _ := r.FormFile("file")

	// you can check here

	// ensure the file name is unique
	fileName := header.Filename
	absPath := fmt.Sprintf("%v/%v", rootPath, fileName)

	// if the file exist, open the file. Otherwise, create a new file
	local, err := os.Open(absPath)
	defer local.Close() // do not forget to close the file
	if err != nil {
		local, _ = os.Create(absPath)
	}

	// create a buffer to read the file
	buffer := make([]byte, 1024)

	// record the read size of the file
	total, err := rdb.Get(ctx, fmt.Sprintf("%v%v", HistoryKey, fileName)).Int()
	// the file have not been uploaded yet
	if err == nil {
		// offset proper size
		remote.Seek(int64(total), io.SeekStart)
	}

	for {
		// loop to read the file
		size, err := remote.Read(buffer)

		// EOF => finish reading
		if err == io.EOF {
			break
		}

		// save buffer to the local file
		local.WriteAt(buffer, int64(total))

		// Add total size
		total = total + size

		// save total size to the redis
		if err = rdb.Set(ctx, fmt.Sprintf("%v%v", HistoryKey, fileName), total, time.Minute*10).Err(); err != nil {
			log.Println(err.Error())
		}
	}

	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":9090", nil)
}
