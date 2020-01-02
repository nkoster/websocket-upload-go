# Websocket Upload Golang

Upload a huge file in chunks over a websocket connection.

Based on [websocket-upload](https://github.com/nkoster/websocket-upload), where the server is written in JavaScript.

I'm using the Golang [Gorilla](https://github.com/gorilla/websocket/) websocket framework for the server. 
Front-end JavaScript is based on a gist from Alessandro Diaferia.
https://gist.github.com/alediaferia/cfb3a7503039f9278381

DISCLAIMER: I'm not sure if this is a good idea, but it actually works.
This is currently an experiment in progress. I'm very open for comments. Also, this is my first Golang thingy.

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
An uploaded file will appear in /tmp/, in this example.

You can adjust the host and the port:

```
./websocket-upload -host example.com -port 8000
```

Before uploading, an MD5 sum is calculated. This is for the future.
I'm using [js-spark-md5](https://github.com/satazor/js-spark-md5) from André Cruz
for the in-fly MD5 calculation.
