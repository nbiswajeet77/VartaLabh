package main

import (
	"log"
	"net/http"

	"vartalabh.com/m/accounts"
	"vartalabh.com/m/agents"
)

func handleRequests() {
	http.HandleFunc("/chat/new", agents.CreateChat)
	http.HandleFunc("/register", accounts.RegisterHandler)
	http.HandleFunc("/login", accounts.LoginHandler)
	http.HandleFunc("/chat/history", agents.GetChatHistory)
	http.HandleFunc("/getChat", agents.GetChat)
	http.HandleFunc("/sendMessage", agents.SendMessage)
	http.HandleFunc("/deleteChat", agents.DeleteChat)
	http.HandleFunc("/editPrompt", agents.EditPrompt)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	agents.DbConn()
	handleRequests()
}
