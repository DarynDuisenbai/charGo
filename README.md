# charGo
 # Simple Chat Application in Golang

This is a simple chat application written in Golang, utilizing WebSocket for real-time communication and PostgreSQL for message storage.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)


## Features

- Real-time chat using WebSocket
- Message storage in PostgreSQL
- User-friendly interface

## Prerequisites

Before you begin, ensure you have met the following requirements:

- [Go (Golang)](https://golang.org/) installed on your machine
- [PostgreSQL](https://www.postgresql.org/) database available
- WebSocket-capable browser (modern browsers)

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/chat-app-golang.git
    ```

2. **Navigate to the project directory:**

    ```bash
    cd chat-app-golang
    ```

3. **Install dependencies:**

    ```bash
    go get -u github.com/gorilla/websocket
    go get -u github.com/joho/godotenv
    go get -u github.com/lib/pq
    ```

4. **Create a `.env` file in the project root and configure your PostgreSQL connection:**

    ```env
    DB_HOST=your_database_host
    DB_PORT=your_database_port
    DB_USER=your_database_user
    DB_PASSWORD=your_database_password
    DB_NAME=your_database_name
    ```

5. **Build the application:**

    ```bash
    go build
    ```

## Usage

1. **Start your PostgreSQL database.**

2. **Run the application:**

    ```bash
    ./chat-app-golang
    ```

3. **Open your web browser and navigate to [http://localhost:8080](http://localhost:8080).**

4. **Enjoy chatting with your friends!**

##Creators
Asylzhan, Daryn Se-2214



