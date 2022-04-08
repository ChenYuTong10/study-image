# Example

## Usage

The example requires a working Go development environment.

To use the example, run the `main.go` in the command line and open `index.html`.

## Frontend

The frontend code is in [index.html](https://github.com/ChenYuTong10/study-image/blob/master/second/index.html).

Generally, the whole flow is roughly same as the image upload.

Firstly, calculating the `md5` of the image.
In the example, we use the `spark-md5` package to calculate the md5.

```JavaScript

function GetImageMd5(image) {
    return new Promise(resolve => {
        // read image
        const reader = new FileReader();
        reader.readAsBinaryString(image);
        reader.onload = (e) => {
            // calculate the md5 of the image
            const spark = new SparkMD5();
            spark.appendBinary(e.target.result);
            resolve(spark.end());
        };
    });
}

```

Next, we send a check request instead of sending the upload request to verify whether the image has been uploaded.

```JavaScript

function CheckImage(md5) {
    const form = new FormData();
    form.append("md5", md5);

    return axios.post(`${url}check`, form);
}

``` 

Then, according to the response, we decide to upload or not upload.

```JavaScript

// check the image whether it is exist on server
let result = await CheckImage(md5);

if(result.data === "no") {
    // no same image on server, so we need to upload image
    result = await UploadImage(md5, image);

    // display the image
    if(result.status === 200) {
        img.src = `${url}${result.data}`;
    }

    return
}

// if the image has exist, diplay the image directly
img.src = `${url}${result.data}`;

```

If the image is exist on the server, we only send a check request instead of a heavy upload request.
That is why the image can be upload in a second.

## Backend

The backend code is in [main.go](https://github.com/ChenYuTong10/study-image/blob/master/second/main.go).

Corresponding to the frontend, the backend offers two interfaces.

One to check the md5 of the image and another one to accept the upload image.

To get more details, you can see the `main.go`.

## Work flow

![alt work-flow](https://github.com/ChenYuTong10/study-image/blob/master/second/flow.png)