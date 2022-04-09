# Example

## Usage

The example requires a working Go development environment.

To use the example, run the `main.go` in the command line and open `index.html`.

## Frontend

The core to frontend is how to slice the image. Frontend can slice the image in equal step by using the `slice` function.

After slicing, frontend sends upload request with the chunk and its unique sign such as using the combination of timestamp and userId.

Finally, sending request to tell the server merge the image if all the requests succeed. Using `Promise.all` function can easily to check the status of all upload requests.

You can see more details in the [index.html](https://github.com/ChenYuTong10/study-image/blob/master/slice/index.html).

## Backend

Backend need to receive the chunk and use symbol to sign order of the chunk. The symbol usually is from the frontend.At the same time, we make a directory to store these chunks until the merge request comes.

Pay attention to the chunk and it is a **file** not byte flow.

When then merge request comes, create a new file to store the all the contents. Then, tranverse all the chunks blew the appointed directory, read it and write to the new file created above.

You can see more details in the [main.go](https://github.com/ChenYuTong10/study-image/blob/master/slice/main.go).

## Work flow

![alt work-flow](https://github.com/ChenYuTong10/study-image/blob/master/slice/flow.png)