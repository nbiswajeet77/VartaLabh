package main

import (
	"log"
	"net/http"

	"vartalabh.com/m/agents"
)

func handleRequests() {
	http.HandleFunc("/chat", agents.Chat)

	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
