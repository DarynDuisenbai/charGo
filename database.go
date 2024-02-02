package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// ChatDB представляет собой класс для работы с базой данных.
type ChatDB struct {
	db *sql.DB
}

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

// Message представляет собой модель сообщения
type Message struct {
	ID       int
	Content  string
	Username string
}

// NewChatDB создает новый экземпляр ChatDB с подключением к базе данных.
func NewChatDB() (*ChatDB, error) {
	db, err := sql.Open("postgres", "user=postgres password=mother1978 dbname=chat sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Создаем таблицу messages, если она еще не существует
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS messages (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50),
        content TEXT,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
`)
	if err != nil {
		log.Fatal("Error creating messages table:", err)
		return nil, err
	}

	return &ChatDB{db: db}, nil
}

func (c *ChatDB) CreateUser(username, passwordHash string) error {
	_, err := c.db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, passwordHash)
	return err
}

func (c *ChatDB) GetUserByName(username string) (*User, error) {
	var user User
	err := c.db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	return &user, err
}

// Close закрывает соединение с базой данных.
func (c *ChatDB) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c *ChatDB) GetMessages(filterUsername string, sortByTimestamp bool) ([]string, error) {
	var rows *sql.Rows
	var err error

	// Используйте параметры filterUsername и sortByTimestamp для формирования SQL-запроса
	if filterUsername != "" {
		rows, err = c.db.Query("SELECT content FROM messages WHERE username = $1 ORDER BY timestamp ASC", filterUsername)
	} else {
		if sortByTimestamp {
			rows, err = c.db.Query("SELECT content FROM messages ORDER BY timestamp ASC")
		} else {
			rows, err = c.db.Query("SELECT content FROM messages")
		}
	}

	if err != nil {
		log.Println("Error querying messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			log.Println("Error scanning message:", err)
			return nil, err
		}
		messages = append(messages, content)
	}

	return messages, nil
}

// InsertMessage вставляет сообщение в базу данных.
func (c *ChatDB) InsertMessage(userID string, content string) error {
	_, err := c.db.Exec("INSERT INTO messages (user_id, content) VALUES ($1, $2)", userID, content)
	return err
}

// Добавьте новый метод для получения сообщений с временной меткой
func (c *ChatDB) GetMessagesWithTimestamp(filterUsername string, sortByTimestamp bool) ([]string, error) {
	var rows *sql.Rows
	var err error

	// Используйте параметры filterUsername и sortByTimestamp для формирования SQL-запроса
	if filterUsername != "" {
		rows, err = c.db.Query("SELECT content, timestamp FROM messages WHERE username = $1 ORDER BY timestamp ASC", filterUsername)
	} else {
		if sortByTimestamp {
			rows, err = c.db.Query("SELECT content, timestamp FROM messages ORDER BY timestamp ASC")
		} else {
			rows, err = c.db.Query("SELECT content, timestamp FROM messages")
		}
	}

	if err != nil {
		log.Println("Error querying messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var content string
		var timestamp time.Time
		if err := rows.Scan(&content, &timestamp); err != nil {
			log.Println("Error scanning message:", err)
			return nil, err
		}
		formattedMsg := fmt.Sprintf("[%s] %s", timestamp.Format("2006-01-02 15:04:05"), content)
		messages = append(messages, formattedMsg)
	}

	return messages, nil
}
