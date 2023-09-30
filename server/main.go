package main

import (
	"log"
	"net/http"

	"vartalabh.com/m/accounts"
	"vartalabh.com/m/agents"
)

func handleRequests() {
	http.HandleFunc("/chat", agents.Chat)
	http.HandleFunc("/register", accounts.RegisterHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
