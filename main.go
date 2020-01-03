package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	var serverHost string = "localhost"
	var serverPort string = "8086"
	var serverStore string = "/tmp"
	for i, arg := range os.Args {
		if arg == "-host" {
			if i+1 < len(os.Args) {
				serverHost = os.Args[i+1]
			}
			arg = ""
		}
		if arg == "-port" {
			if i+1 < len(os.Args) {
				serverPort = os.Args[i+1]
			}
			arg = ""
		}
		if arg == "-store" {
			if i+1 < len(os.Args) {
				serverStore = os.Args[i+1]
				serverStore = strings.TrimRight(serverStore, "/")
			}
			arg = ""
		}
		if arg == "--help" {
			log.Printf("\nusage: %s [[-host <host>] [-port <port>]]\n", os.Args[0])
			os.Exit(0)
		}
	}
	log.Println("flextube", serverHost+":"+serverPort, serverStore)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html(serverHost, serverPort))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var err error
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("error", err)
			return
		}
		log.Println("connection", r.RemoteAddr)
		filename := ""
		var f *os.File
		go func(conn *websocket.Conn) {
			for {
				mt, data, connErr := conn.ReadMessage()
				if connErr != nil {
					log.Println("error", connErr)
					return
				}
				if mt == 1 {
					event := strings.Split(string(data), ":")
					if event[0] == "upload" {
						filename = serverStore + "/" + event[1]
						f, err = os.Create(filename)
						if err != nil {
							log.Println(err)
						}
					}
					log.Println(string(event[0]), filename)
					if event[0] == "ready" {
						filename = ""
						f.Close()
						if err := conn.WriteMessage(1, []byte("ready")); err != nil {
							log.Println("error sending ready message")
						}
					}
				}
				if mt == 2 {
					log.Println("chunk", filename)
					f.Write(data)
				}
			}
		}(conn)
	})
	log.Fatal(http.ListenAndServe(serverHost+":"+serverPort, nil))
}
