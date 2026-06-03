package models

import (
	"fmt"
	"log"
	"os"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func InitDB() {
	if os.Getenv("ENV") == "" {
		_ = godotenv.Load("../.env")
	}
	env := os.Getenv("ENV")
	var user, password, host, dbname, port string

	switch env {
	case "DEV":
		user = os.Getenv("DEV_DB_USER")
		password = os.Getenv("DEV_DB_PASSWORD")
		host = os.Getenv("DEV_DB_HOST")
		dbname = os.Getenv("DEV_DB_NAME")
		port = os.Getenv("DEV_DB_PORT")
	case "TEST":
		user = os.Getenv("TEST_DB_USER")
		password = os.Getenv("TEST_DB_PASSWORD")
		host = os.Getenv("TEST_DB_HOST")
		dbname = os.Getenv("TEST_DB_NAME")
		port = os.Getenv("TEST_DB_PORT")
	case "PROD":
		user = os.Getenv("PROD_DB_USER")
		password = os.Getenv("PROD_DB_PASSWORD")
		host = os.Getenv("PROD_DB_HOST")
		dbname = os.Getenv("PROD_DB_NAME")
		port = os.Getenv("PROD_DB_PORT")
	default:
		log.Fatalf("Unknown ENV: %s", env)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect database: %v", err)
	}

	DB = db

	// Auto-migrate schema
	if err := db.AutoMigrate(&User{}, &Task{}); err != nil {
		log.Printf("Migration warning: %v", err)
	}
}

func SeedTestData(db *gorm.DB){
	env := os.Getenv("ENV")
	if env == "TEST"{
		if err := db.Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE;").Error; err != nil {
			log.Fatalf("Failed to reset task table: %v", err)
		}

		if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE;").Error; err != nil {
			log.Fatalf("Failed to reset user table: %v", err)
		}

		user1 := User{Email: "leon@gmail.com"}
		user2 := User{Email: "youssef@hotmail.com"}
		
		var err error
		user1.Password, err = HashPassword("leon123")
		
		if err != nil{
			log.Printf("Error hashing password: %v", err)
		}
		user2.Password, err = HashPassword("youssef123")

		if err != nil{
			log.Printf("Error hashing password: %v", err)
		}

		if err := db.Create(&[]User{user1, user2,}).Error; err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}

		// Seed default tasks
		if err := db.Create(&[]Task{
			{Title: "Learn Go", Description: "Study Go basics", Completed: false, UserID: 1},
			{Title: "Build API", Description: "Create a REST API", Completed: false, UserID: 2},
		}).Error; err != nil{
			log.Fatalf("Failed to seed tasks: %v", err)
		}

	}

}

func HashPassword(password string) (string, error) {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}
