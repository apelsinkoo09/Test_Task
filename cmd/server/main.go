package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"test_task/internal/handler"

	_ "github.com/lib/pq"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var user Config
	err = json.Unmarshal(file, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DBConnect(c Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname='%s' host=%s sslmode=%s port=%d",
		c.User,
		c.Password,
		c.DBName,
		c.Host,
		c.SSLMode,
		c.Port,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return db, nil
}

func main() {
	config, err := LoadConfig("configs/config.json")
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
	}
	db, err := DBConnect(*config)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	http.HandleFunc("/api/v1/wallet", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}
		handler.UpdateBalanceHandler(db)(w, r)
	})
	http.HandleFunc("/api/v1/wallets/", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}
		handler.GetUUIDBalanceHandler(db)(w, r)
	})

	// Запуск сервера
	log.Printf("Server is running on 8081")
	serv := http.ListenAndServe(":8081", nil)
	if serv != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
