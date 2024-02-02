package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	clients     map[*client]bool
	join        chan *client
	leave       chan *client
	forward     chan message
	db          *ChatDB
	userCounter int
}

type message struct {
	client  *client
	content string
}

func newRoom(db *ChatDB) *room {
	return &room{
		forward: make(chan message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		db:      db,
	}
}

func (r *room) assignUsername() string {
	r.userCounter++
	return fmt.Sprintf("user%d", r.userCounter)
}

func (r *room) run(filterUsername string, sortByTimestamp bool) {
	for {
		select {
		case client := <-r.join:
			username := r.assignUsername()
			r.clients[client] = true
			client.username = username

			// Отправим существующие сообщения новому клиенту
			r.sendExistingMessages(client, filterUsername, sortByTimestamp)

			// Отправим уведомление о новом пользователе всем клиентам
			r.sendSystemMessage(fmt.Sprintf("%s joined the chat.", username))

		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)

			// Отправим уведомление об отключившемся пользователе всем клиентам
			r.sendSystemMessage(fmt.Sprintf("%s left the chat.", client.username))

		case msg := <-r.forward:
			// Добавляем сохранение сообщения в базу данных
			err := r.db.InsertMessage(msg.client.username, msg.content)
			if err != nil {
				log.Println("Error inserting message into database:", err)
			}

			// Обновляем код для отправки сообщений всем клиентам
			for otherClient := range r.clients {
				otherClient.receive <- []byte(msg.content)
			}

		}
	}
}

// Добавьте новый метод для отправки существующих сообщений клиенту с временной меткой
func (r *room) sendExistingMessages(client *client, filterUsername string, sortByTimestamp bool) {
	// Получаем сообщения из базы данных
	messages, err := r.db.GetMessagesWithTimestamp(filterUsername, sortByTimestamp)
	if err != nil {
		log.Println("Error getting messages from database:", err)
		return
	}

	// Отправляем сообщения клиенту
	for _, msg := range messages {
		client.receive <- []byte(msg)
	}
}

func (r *room) sendSystemMessage(message string) {
	for client := range r.clients {
		// Используйте префикс [System] для системных сообщений
		client.receive <- []byte(fmt.Sprintf("[System] %s", message))
	}
}

/*func (r *room) handleRegister(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	// TODO: Добавить валидацию данных

	// Создать хэш пароля
	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Вставить нового пользователя в базу данных
	err = r.db.CreateUser(username, hashedPassword)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *room) handleLogin(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	// Получить пользователя из базы данных по имени
	user, err := r.db.GetUserByName(username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Проверить пароль
	if !checkPassword(password, user.PasswordHash) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// TODO: Реализовать механизм аутентификации и хранения сессии (пример: использование cookie)
	// Здесь можно использовать стандартную библиотеку net/http для управления сессиями.

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

*/

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

/*func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// checkPassword проверяет соответствие пароля хешу
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

*/
