package main

import (
	"fmt"
	"os"
	"strings"
)

func args() (string, string, string) {
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
		if arg == "--help" || arg == "-help" || arg == "/h" {
			fmt.Printf("usage: %s [[-host <host>] [-port <port>] [-store <path>]]\n", os.Args[0])
			fmt.Println("note: -store <path> must be absolute (/)")
			os.Exit(0)
		}
	}
	return serverHost, serverPort, serverStore
}
