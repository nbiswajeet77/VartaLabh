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
	http.HandleFunc("/login", accounts.LoginHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
