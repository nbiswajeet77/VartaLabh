package accounts

import (
	"bufio"
	"fmt"
	"log"
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
	if r.Method == "POST" {
		db := agents.DbConn()
		fmt.Printf("Enter First Name: ")
		firstName := takeInput()
		fmt.Printf("Enter Last Name: ")
		lastName := takeInput()
		fmt.Printf("Enter Password: ")
		pass := takeInput()
		fmt.Printf("%s, %s\n", firstName, lastName)

		password, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}

		_, err = db.Exec("INSERT INTO user(firstname, lastname,password) VALUES(?,?,?)", firstName, lastName, password)
		if err != nil {
			fmt.Println("Error when inserting: ", err.Error())
			panic(err.Error())
		}
		log.Println("=> Inserted: First Name: " + firstName + " | Last Name: " + lastName)
	}
}
