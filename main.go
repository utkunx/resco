package main

import (
	"log"
	"net/http"
	"os"
	"resco/db"
	"resco/handlers"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables or defaults")
	}

	// Load database configuration from environment variables
	dbConfig := db.Config{
		Server:   getEnv("DB_SERVER", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 1433),
		User:     getEnv("DB_USER", "sa"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_DATABASE", "RESCO_2019"),
	}

	// Initialize database connection
	err = db.InitDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	log.Println("Database connected successfully")

	// Create router
	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	router.HandleFunc("/api/bom/{itemCode}", handlers.GetBOMByItemCode).Methods("GET")
	router.HandleFunc("/api/bomcn/{itemCode}", handlers.GetBOMByItemCodeCN).Methods("GET")
	router.HandleFunc("/api/bomcombined/{itemCode}", handlers.GetBOMByItemCodeCombined).Methods("GET")
	router.HandleFunc("/api/bomtotal/{itemCode}", handlers.GetBOMTotal).Methods("GET")
	router.HandleFunc("/api/queryhe/{itemCode}", handlers.QueryHeihu).Methods("GET")
	router.HandleFunc("/api/checkproduct/{itemCode}", handlers.CheckProduct).Methods("GET")

	// Get server port from environment or use default
	port := getEnv("PORT", "8080")

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get environment variable as int with default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}