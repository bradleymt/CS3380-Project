package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func migrateTables() error {
	err := DB.AutoMigrate(
		&User{},
		&Publisher{},
		&Family{},
		&Game{},
		&Discount{},
		&Purchase{},
		&Refund{},
		&Library{},
		&LibraryItem{},
		&Wishlist{},
		&WishlistItem{},
		&Cart{},
		&CartItem{},
	)

	return err
}

func ConnectDatabase() error {

	// Retrieve environment variables with defaults
	user := GetEnv("POSTGRES_USER", "user")
	password := GetEnv("POSTGRES_PASSWORD", "password")
	dbName := GetEnv("POSTGRES_DB", "db")
	dbHost := GetEnv("POSTGRES_HOST", "postgres")
	dbPort := GetEnv("POSTGRES_PORT", "5432")

	// Build the connection string
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, user, password, dbName,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond,                         // Slow SQL threshold
			LogLevel:                  logger.Info | logger.Error | logger.Warn, // Log level
			IgnoreRecordNotFoundError: false,                                    // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,                                     // Don't include params in the SQL log
			Colorful:                  true,                                     // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	DB = db
	return migrateTables()
}
