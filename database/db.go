package database

import (
	"fmt"
	"log"
	"os"

	"github.com/untrik/FromSkateToZOH/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbNAME := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s connect_timeout=60 dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbPort, dbNAME,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Проблема с подключением к базе данных:", err)
	}
	log.Println("Подключение к базе данных установлено.")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Student{},
		&models.Product{},
		&models.Event{},
		&models.EventParticipant{},
		&models.Reward{},
		&models.Order{},
	)
	if err != nil {
		log.Fatal("Ошибка миграции таблиц:", err)
	}
}
