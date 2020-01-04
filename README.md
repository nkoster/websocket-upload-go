# Websocket Upload Golang

Upload a huge file in chunks over a websocket connection.

Based on [websocket-upload](https://github.com/nkoster/websocket-upload), where the server is in JavaScript / nodejs.

I'm using the Golang [Gorilla](https://github.com/gorilla/websocket/) websocket framework for the server. 
Front-end JavaScript is based on a gist from Alessandro Diaferia.
https://gist.github.com/alediaferia/cfb3a7503039f9278381

DISCLAIMER: I'm not sure if this is a good idea, but it actually works.
This is currently an experiment in progress. I'm very open for comments. Also, this is my first Golang experience.

Usage, assuming you have your Go environment prepared:

```
git clone https://github.com/nkoster/websocket-upload-go
cd websocket-upload-go
go get github.com/gorilla/websocket
go build
./websocket-upload-go
````

or

```
go run *.go
```

Open http://localhost:8086 and drag-and-drop a file in the page.
An uploaded file will appear in /tmp/, but you can adjust that:

```
./websocket-upload-go -store /store
```

The ```-store``` path must be absolute.

Before uploading, an MD5 sum is calculated. 
The MD5sum will be used as file name, and the original file name will be saved as a symlink:

```
/store/
  ├── files
  │   └── 166c5a55e29a73db2afd997b52e6e554
  └── links
      └── my-video.mp4 -> /store/files/166c5a55e29a73db2afd997b52e6e554
 ```

I'm using [js-spark-md5](https://github.com/satazor/js-spark-md5) from André Cruz
for the incremental (stream) MD5 calculation.

You can adjust the host and the port:

```
./websocket-upload-go -host example.com -port 8000
```
