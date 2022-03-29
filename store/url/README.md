# Example

## Usage

The example requires a working Go development environment.

To use the example, run the `main.go` in the command line and open `index.html`.

## Database

The following code is the ddl of the database.

```SQL
/* mysql  Ver 8.0.27 for Linux on x86_64 (MySQL Community Server - GPL) */

CREATE TABLE t_url_image
(
    id INTEGER auto_increment,
    path VARCHAR(100),
    PRIMARY KEY (id)
);
```

## Frontend

The frontend code is in [index.html](https://github.com/ChenYuTong10/study-image/blob/master/store/url/index.html).

When the user upload an image, the **change** event will be triggered. It will
get the upload file and send request directly.

When you want to show the image, you need to get the path of it on the server.
In the example, when the server handle image successfully, it will give back the path.
So you can access the image through the path.

## Backend

The backend code is in [main.go](https://github.com/ChenYuTong10/study-image/blob/master/store/url/main.go)

As you see, the backend use the file server to realize the access of image. You can also use the 
function `http.ServeFile`.

```Golang
func getImage(w http.ResponseWriter, r *http.Request) {

    // get the image name
    name := request.PostFormValue("name")

    // get the root path
    rootPath, _ := os.Getwd()

    // check the existence of the image
    absPath := fmt.Sprintf("%v/%v", rootPath, name)
    _, err := os.Stat(absPath)
    if os.IsNotExist(err) {
        w.Write([]byte("not exist"))
        return
    }

    // serve the image
    http.ServeFile(w, r, absPath)
}
```

## Hint

**Be careful with the path diffence between the Linux os and Windows os.**
The image on the Linux os is separated with '/', but it is '\\' on the Windows os.