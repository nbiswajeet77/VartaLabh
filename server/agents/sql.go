package agents

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"vartalabh.com/m/model"
)

var db *sql.DB

func DbConn() {
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
}

func FetchParticularChat(chatID string) (*model.GetChatResponse, error) {
	var resp *model.GetChatResponse
	checkChat, err := db.Query("SELECT chatID, prompt, messages FROM Chats WHERE chatID=?", chatID)
	if err != nil {
		return nil, err
	}
	defer checkChat.Close()
	for checkChat.Next() {
		var chatID, prompt string
		var msg []byte
		err = checkChat.Scan(&chatID, &prompt, &msg)
		if err != nil {
			return nil, err
		}
		var messages []model.Message
		if err := json.Unmarshal(msg, &messages); err != nil {
			return nil, err
		}
		resp = &model.GetChatResponse{
			ChatId:   chatID,
			Prompt:   prompt,
			Messages: messages,
		}
	}
	return resp, nil
}

func DeleteParticularChat(chatID string) (int64, error) {
	query := "DELETE FROM chats WHERE chatId = ?"

	result, err := db.Exec(query, chatID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func FetchUserChats(userID string) ([]*model.ChatHistoryResponse, error) {
	chathistory := make([]*model.ChatHistoryResponse, 0)
	checkChats, err := db.Query("SELECT chatID, prompt FROM Chats WHERE userID=?", userID)
	if err != nil {
		return nil, err
	}
	defer checkChats.Close()
	for checkChats.Next() {
		var chatId, prompt string
		err = checkChats.Scan(&chatId, &prompt)
		if err != nil {
			return nil, err
		}
		chathistory = append(chathistory, &model.ChatHistoryResponse{
			ChatId: chatId,
			Prompt: prompt,
		})
	}
	return chathistory, nil
}

func FetchUser(userID string) *model.User {
	user := &model.User{}
	checkUser, err := db.Query("SELECT userID, password FROM Users WHERE userID=?", userID)
	if err != nil {
		panic(err.Error())
	}
	defer checkUser.Close()
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
	_, err := db.Exec("INSERT INTO Users(userID,password) VALUES(?,?)", userID, password)
	if err != nil {
		fmt.Println("Error when inserting in users table: ", err.Error())
		return err
	}
	return nil
}

func CreateChatEntry(userID, chatID, prompt string, messages []byte) error {
	_, err := db.Exec("INSERT INTO Chats(userID,chatID,messages,prompt) VALUES(?,?,?,?)", userID, chatID, messages, prompt)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChatEntry(chatId, prompt string, messages []byte) error {
	_, err := db.Exec("UPDATE chats SET prompt = ?, messages = ? WHERE chatId = ?", prompt, messages, chatId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserCurrentChat(userID, chatId string) error {
	_, err := db.Exec("UPDATE Users SET chatId = ? WHERE userId = ?;", chatId, userID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChatMessages(chatId string, messages []byte) error {
	_, err := db.Exec("UPDATE Chats SET messages = ? WHERE chatID = ?;", messages, chatId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserPrompt(userID, prompt string) {
	_, err := db.Exec("UPDATE Users SET prompt = ? WHERE email = ?;", prompt, userID)
	if err != nil {
		fmt.Println("Error when Updating for user ID: ", userID, err.Error())
		panic(err.Error())
	}
	log.Println("=> Updated: UserID: " + userID)
}
