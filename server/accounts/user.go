package accounts

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"vartalabh.com/m/agents"
)

func takeInput() string {
	in := bufio.NewReader(os.Stdin)
	input, _ := in.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if true { //r.Method == "POST" {
		fmt.Printf("Enter Email ID: ")
		emailID := takeInput()
		fmt.Printf("Enter Password: ")
		pass := takeInput()
		fmt.Printf("%s, %s\n", emailID, pass)
		prompt := "you are a mental health counsellor. Talk to user, ask repititive questions,keep the conversation going. Also ask the user how much progress he made based on the prompt provided."

		password, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}

		agents.CreateUser(emailID, prompt, password)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if true {
		fmt.Printf("Enter Email ID: ")
		emailID := takeInput()
		fmt.Printf("Enter Password: ")
		pass := takeInput()
		fmt.Printf("%s, %s\n", emailID, pass)

		if strings.Trim(emailID, " ") == "" || strings.Trim(pass, " ") == "" {
			fmt.Println("Parameter's can't be empty")
			return
		}

		user := agents.FetchUser(emailID)

		errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
		if errf != nil { //Password does not match!
			fmt.Println(errf)
		} else {
			fmt.Println("User Logged in successfully")
		}
	}
}
