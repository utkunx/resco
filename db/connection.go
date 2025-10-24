package db

import (
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
)

type Config struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
}

var DB *sql.DB

// InitDB initializes the database connection
func InitDB(config Config) error {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		config.Server, config.User, config.Password, config.Port, config.Database)

	var err error
	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}