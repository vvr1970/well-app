package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// InitDB инициализирует подключение к PostgreSQL
func InitDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Создание таблицы, если её нет
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS wells (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			location VARCHAR(100) NOT NULL,
			gammag DECIMAL(6,4),
			temp DECIMAL(5,2),
			tempust DECIMAL(5,2),
			depth DECIMAL(6,1),
			pbuf DECIMAL(6,3),
			ptb DECIMAL(6,3),
			ppl DECIMAL(6,3),
			pz DECIMAL(6,3),
			q DECIMAL(6,3),
			roughness DECIMAL(6,5),
			diametr DECIMAL(4,1),
			a DECIMAL(17,16),
			b DECIMAL(17,16),
			mu DECIMAL(3,2),
			wgf DECIMAL(3,1),
			rog DECIMAL(5,1),
			hw DECIMAL(5,1),
			qmin DECIMAL(6,3),
			pmax DECIMAL(6,3),
			status VARCHAR(100)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}
