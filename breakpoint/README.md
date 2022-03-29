# Example

## Usage

The example requires a working Go development environment.

To use the example, run the `main.go` in the command line and open `index.html`.

## Frontend

The frontend code is in [index.html](https://github.com/ChenYuTong10/study-image/blob/master/breakpoint/index.html).

It is the same with the way which is mentioned in the store directory.

## Backend

The backend code is in [main.go](https://github.com/ChenYuTong10/study-image/blob/master/breakpoint/main.go).

The core to realize the breakpoint resume is **save the history upload size of the file**.
You can use anything to store the size like a text file or a redis database.

First, you need to make a buffer to read the file.

```Golang
    buffer := make([]byte, 1024)
``` 

Next, prepare a counter to record the size which has been read.

Then, loop to read the file and add the read size to the counter until the file has been read completely.

```Golang
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

        // do something ...
    }
```

But when to save the read size?

When the upload process was interrupted, it is too late to save the read size. Therefore, we will save the size in every loop.

```Golang
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

        // save the read size...
        // you can use any way to save it
        if err = rdb.Set(ctx, fmt.Sprintf("%v%v", HistoryKey, fileName), total, time.Minute*10).Err(); err != nil {
            log.Println(err.Error())
        }
    }
```

What can we do with the stored size?

The next time when the same file is uploaded, we can get its history upload size through the name or anything else.
Then we can skip the stored size and read behind directly.

```Golang
    // record the read size of the file
    total, err := rdb.Get(ctx, fmt.Sprintf("%v%v", HistoryKey, fileName)).Int()
    // the file have not been uploaded yet
    if err == nil {
        // offset proper size
        remote.Seek(int64(total), io.SeekStart)
    }
``` 

Above all, the breakpoint resume will realize through following steps:

01. get upload file
02. Open or create the file
03. make a proper buffer
04. get the upload size of the file
05. loop to read and add read size to the counter
06. close the file
