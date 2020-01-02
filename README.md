# Web Socket Upload Golang

Upload a huge file in chunks over a web socket connection.

(Based on https://github.com/nkoster/websocket-upload, where the server is written in JavaScript)

I'm using the [https://github.com/gorilla/websocket/](Gorilla) web socket framework for the server side web socket logic. 
Front JS is based on a gist from Alessandro Diaferia. (Thank you!)
https://gist.github.com/alediaferia/cfb3a7503039f9278381

DISCLAIMER: I'm not sure if this is a good idea, but it actually works.
This is currently an experiment in progress. I'm very open for comments.

Usage, assuming you have your Go environment up & running:

```
git clone https://github.com/nkoster/websocket-upload-go
cd websocket-upload-go
go build
./websocket-upload-go
````

Open http://localhost:8086 and drag-and-drop a file in the page.
An uploaded file will appear in /tmp/, in this example.

Before uploading, an MD5 sum is calculated. This is for the future.
I'm using [https://github.com/satazor/js-spark-md5](js-spark-md5) from André Cruz
for the in-fly MD5 calculation. Obrigado!
