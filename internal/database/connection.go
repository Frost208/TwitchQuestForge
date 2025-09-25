package database

import (
        "context"
        "database/sql"
        "fmt"
        "os"
        "time"
        
        _ "github.com/go-sql-driver/mysql"
)

// DB represents the database connection
var DB *sql.DB

// Connect establishes a connection to the MySQL database
func Connect() error {
        // Get database configuration from environment
        dbHost := os.Getenv("DB_HOST")
        dbPort := os.Getenv("DB_PORT")
        dbUser := os.Getenv("DB_USER")
        dbPassword := os.Getenv("DB_PASSWORD")
        dbName := os.Getenv("DB_NAME")
        
        if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
                return fmt.Errorf("missing required database environment variables: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
        }
        
        // Create connection string with timeout
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
                dbUser, dbPassword, dbHost, dbPort, dbName)
        
        // Open database connection
        db, err := sql.Open("mysql", dsn)
        if err != nil {
                return fmt.Errorf("failed to open database: %v", err)
        }
        
        // Test the connection with context timeout
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        if err := db.PingContext(ctx); err != nil {
                db.Close()
                return fmt.Errorf("failed to ping database: %v", err)
        }
        
        // Configure connection pool
        db.SetMaxOpenConns(25)
        db.SetMaxIdleConns(5)
        
        DB = db
        fmt.Println("Successfully connected to database")
        return nil
}

// Close closes the database connection
func Close() error {
        if DB != nil {
                return DB.Close()
        }
        return nil
}