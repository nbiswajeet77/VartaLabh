package agents

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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

func FetchParticularChat(chatID string) (*model.Chat, error) {
	var resp *model.Chat
	checkChat, err := db.Query("SELECT chatID, prompt, messages, summary FROM Chats WHERE chatID=?", chatID)
	if err != nil {
		return nil, err
	}
	defer checkChat.Close()
	for checkChat.Next() {
		var chatID, prompt, summary string
		var msg []byte
		err = checkChat.Scan(&chatID, &prompt, &msg, &summary)
		if err != nil {
			return nil, err
		}
		var messages []model.Message
		if err := json.Unmarshal(msg, &messages); err != nil {
			return nil, err
		}
		resp = &model.Chat{
			ChatId:   chatID,
			Prompt:   prompt,
			Messages: messages,
			Summary:  summary,
		}
	}
	return resp, nil
}

func DeleteParticularChat(chatID string) (int64, error) {
	query := "DELETE FROM Chats WHERE chatID = ?"

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
	checkChats, err := db.Query("SELECT chatID, prompt, summary, messages, createdAt, updatedAt FROM Chats WHERE userID=? ORDER BY updatedAt ASC;", userID)
	if err != nil {
		return nil, err
	}
	defer checkChats.Close()
	for checkChats.Next() {
		var chatId, prompt, summary string
		var createdAt, updatedAt time.Time
		var msg []byte
		err = checkChats.Scan(&chatId, &prompt, &summary, &msg, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		var messages []model.Message
		if err := json.Unmarshal(msg, &messages); err != nil {
			return nil, err
		}
		if len(messages) > 2 {
			chathistory = append(chathistory, &model.ChatHistoryResponse{
				ChatId:    chatId,
				Prompt:    prompt,
				Summary:   summary,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			})
		}
	}
	length := len(chathistory)

	for i := 0; i < length/2; i++ {
		chathistory[i], chathistory[length-1-i] = chathistory[length-1-i], chathistory[i]
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
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec("INSERT INTO Users(userID,password,createdAt,updatedAt) VALUES(?,?,?,?)", userID, password, createdAt, createdAt)
	if err != nil {
		fmt.Println("Error when inserting in users table: ", err.Error())
		return err
	}
	return nil
}

func CreateChatEntry(userID, chatID, prompt string, messages []byte) error {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec("INSERT INTO Chats(userID,chatID,messages,prompt,summary,createdAt,updatedAt) VALUES(?,?,?,?,?,?,?)", userID, chatID, messages, prompt, "", createdAt, createdAt)
	if err != nil {
		return err
	}
	return nil
}

func CreateWaitlistEntry(emailId string) error {
	_, err := db.Exec("INSERT INTO Emails(emailID) VALUES(?)", emailId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChatEntry(chatId, prompt, summary string, messages []byte) error {
	_, err := db.Exec("UPDATE Chats SET prompt = ?, messages = ?, summary = ? WHERE chatID = ?", prompt, messages, summary, chatId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChatSummary(chatId, summary string) error {
	_, err := db.Exec("UPDATE Chats SET summary = ? WHERE chatID = ?", summary, chatId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserCurrentChat(userID, chatId string) error {
	if userID == "" {
		return nil
	}
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	var uid = userID
	if userID == nil {
		uid = uuid.New().String()
	}
	_, err := db.Exec("UPDATE Users SET chatId = ?, updatedAt = ? WHERE userId = ?;", chatId, updatedAt, uid)
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
