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
	log.Println("flextube")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html())
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		log.Printf("Connection, %T\n", conn)
		filename := ""
		var f *os.File
		var err error
		go func(conn *websocket.Conn) {
			for {
				mt, data, connErr := conn.ReadMessage()
				if connErr != nil {
					log.Println("Broken connection.", connErr)
					return
				}
				if mt == 1 {
					event := strings.Split(string(data), ":")
					if event[0] == "upload" {
						filename = "/tmp/" + event[1]
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
							log.Println("Error sending ready message")
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
	log.Fatal(http.ListenAndServe(":8086", nil))
}
