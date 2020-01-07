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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	serverHost, serverPort, serverStore, serverHTML := args()
	showStore, showHTML := "", ""
	if serverHTML != "" {
		serverStore = ""
		showHTML = "static:" + serverHTML
	} else {
		showStore = "store:" + serverStore
	}
	log.Println("flextube ->", serverHost+":"+serverPort, showStore, showHTML)
	if serverHTML != "" {
		fs := http.FileServer(http.Dir(serverHTML))
		http.Handle("/", fs)
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, html(serverHost, serverPort))
		})
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var err error
		// not safe, only for dev:
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("error", err)
			return
		}
		log.Println("connection", r.RemoteAddr)
		filename, linkname := "", ""
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
						filename = serverStore + "/files/" + event[1]
						if fileExists(filename) {
							log.Println(filename + " already exists")
							if err := conn.WriteMessage(1, []byte("exists")); err != nil {
								log.Println("error sending exists message")
							}
						} else {
							f, err = os.Create(filename)
							if err != nil {
								log.Println(err)
							}
						}
						linkname = serverStore + "/links/" + event[2]
						err = os.Symlink(filename, linkname)
						if err != nil {
							log.Println(err)
						}
					}
					log.Println(string(event[0]), filename)
					if event[0] == "ready" {
						f.Close()
						if mt := mimeType(filename); mt != "" {
							log.Println(filename, mt)
						} else {
							log.Println(filename, "unknown file type")
						}
						if err := conn.WriteMessage(1, []byte("ready")); err != nil {
							log.Println("error sending ready message")
						}
						filename = ""
					}
				}
				if mt == 2 {
					log.Println("chunk", filename)
					f.Write(data)
				}
			}
		}(conn)
	})
	if err := os.MkdirAll(serverStore+"/files/", 0755); err != nil {
		// log.Println(err)
	}
	if err := os.MkdirAll(serverStore+"/links/", 0755); err != nil {
		// log.Println(err)
	}
	log.Fatal(http.ListenAndServe(serverHost+":"+serverPort, nil))
}
