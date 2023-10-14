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
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func FetchUser(userID string) *model.User {
	db := DbConn()
	user := &model.User{}
	checkUser, err := db.Query("SELECT userID, password FROM Users WHERE userID=?", userID)
	if err != nil {
		panic(err.Error())
	}
	for checkUser.Next() {
		var userID, password string
		err = checkUser.Scan(&userID, &password)
		if err != nil {
			panic(err.Error())
		}
		user.UserId = userID
		user.Password = password
	}
	return user
}

func CreateUser(userID string, password []byte) error {
	db := DbConn()
	_, err := db.Exec("INSERT INTO Users(userID,password) VALUES(?,?)", userID, password)
	if err != nil {
		fmt.Println("Error when inserting: ", err.Error())
		return err
	}
	return nil
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
