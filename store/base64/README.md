# Example

## Usage

The example requires a working Go development environment.

To use the example, run the `main.go` in the command line and open `index.html`.

## Database

The following code is the ddl of the database.

```SQL
/* mysql  Ver 8.0.27 for Linux on x86_64 (MySQL Community Server - GPL) */

CREATE TABLE base64_store_image
(
    id INTEGER auto_increment,
    name VARCHAR(40), -- image name
    code TEXT, -- image code
    PRIMARY KEY (id)
);
```

## Frontend

The frontend code is in [index.html](https://github.com/ChenYuTong10/study-image/blob/master/store/base64/index.html).

When the user upload an image, the **change** event will be triggered. It will
get the upload file, encode the image and send request.

When you want to show the image, you only need to get the base64 encoding of image.

## Backend

The backend code is in [main.go](https://github.com/ChenYuTong10/study-image/blob/master/store/base64/main.go)

When the upload request come, the backend will get the image name and its code.
Next, both of them will be stored in the database.