package agents

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"vartalabh.com/m/model"
)

func DbConn() (db *sql.DB) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	fmt.Println(dbDriver, dbUser, dbPass, dbName)
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func FetchUser(emailID string) *model.User {
	db := DbConn()
	user := &model.User{}
	checkUser, err := db.Query("SELECT email, password, prompt FROM Users WHERE email=?", emailID)
	if err != nil {
		panic(err.Error())
	}
	for checkUser.Next() {
		var email, password, prompt string
		err = checkUser.Scan(&email, &password, &prompt)
		if err != nil {
			panic(err.Error())
		}
		user.Email = email
		user.Password = password
		user.Prompt = prompt
	}
	return user
}

func CreateUser(emailID, prompt string, password []byte) {
	db := DbConn()
	_, err := db.Exec("INSERT INTO Users(email,password,prompt) VALUES(?,?,?)", emailID, password, prompt)
	if err != nil {
		fmt.Println("Error when inserting: ", err.Error())
		panic(err.Error())
	}
	log.Println("=> Inserted: Email: " + emailID)
}

func UpdateUserPrompt(userID, prompt string) {
	db := DbConn()
	_, err := db.Exec("UPDATE Users SET prompt = ? WHERE email = ?;", prompt, userID)
	if err != nil {
		fmt.Println("Error when Updating for user ID: ", userID, err.Error())
		panic(err.Error())
	}
	log.Println("=> Updated: UserID: " + userID)
}
