package main

import (
	"entrytask/internal/http"
	ziherpc "entrytask/internal/rpc"
	"log"
	"net/http"
)

func main() {
	client, err := ziherpc.Dial("tcp", "localhost:8888")
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	httpserver.SetClient(client)
	http.HandleFunc("/login", httpserver.LoginHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
