package accounts

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"vartalabh.com/m/agents"
	"vartalabh.com/m/model"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var user model.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}

		if strings.Trim(user.Email, " ") == "" || strings.Trim(user.Password, " ") == "" {
			WriteOutput(w, "UserID or password can't be empty", http.StatusConflict, err)
			return
		}

		prompt := "you are a mental health counsellor. Talk to user, ask repititive questions,keep the conversation going. Also ask the user how much progress he made based on the prompt provided."

		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			WriteOutput(w, "Error while generating hashed password", http.StatusConflict, err)
			return
		}

		err = agents.CreateUser(user.Email, prompt, password)
		if err != nil {
			WriteOutput(w, "User already registered on application", http.StatusConflict, err)
			return
		}
		WriteOutput(w, "User Registered on application", http.StatusOK, nil)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var user model.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}

		if strings.Trim(user.Email, " ") == "" || strings.Trim(user.Password, " ") == "" {
			WriteOutput(w, "UserID or password can't be empty", http.StatusConflict, err)
			return
		}

		creds := agents.FetchUser(user.Email)

		errf := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(user.Password))
		if errf != nil {
			WriteOutput(w, "Either of userId or password is not correct", http.StatusConflict, err)
			return
		}
		WriteOutput(w, "User Signed in application", http.StatusOK, nil)
	}
}
