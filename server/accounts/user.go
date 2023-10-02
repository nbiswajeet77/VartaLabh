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
	"vartalabh.com/m/model"
)

func takeInput() string {
	in := bufio.NewReader(os.Stdin)
	input, _ := in.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if true { //r.Method == "POST" {
		db := agents.DbConn()
		fmt.Printf("Enter Email ID: ")
		emailID := takeInput()
		fmt.Printf("Enter Password: ")
		pass := takeInput()
		fmt.Printf("%s, %s\n", emailID, pass)

		password, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}

		_, err = db.Exec("INSERT INTO Users(email,password) VALUES(?,?)", emailID, password)
		if err != nil {
			fmt.Println("Error when inserting: ", err.Error())
			panic(err.Error())
		}
		log.Println("=> Inserted: Email: " + emailID + " | Last Name: " + pass)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if true {
		db := agents.DbConn()
		fmt.Printf("Enter Email ID: ")
		emailID := takeInput()
		fmt.Printf("Enter Password: ")
		pass := takeInput()
		fmt.Printf("%s, %s\n", emailID, pass)

		if strings.Trim(emailID, " ") == "" || strings.Trim(pass, " ") == "" {
			fmt.Println("Parameter's can't be empty")
			return
		}

		checkUser, err := db.Query("SELECT email, password FROM Users WHERE email=?", emailID)
		if err != nil {
			panic(err.Error())
		}
		user := &model.User{}
		for checkUser.Next() {
			var email, password string
			err = checkUser.Scan(&email, &password)
			if err != nil {
				panic(err.Error())
			}
			user.Email = email
			user.Password = password
		}
		errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
		if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
			fmt.Println(errf)
		} else {
			fmt.Println("User Logged in successfully")
		}
	}
}
