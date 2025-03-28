package database

import (
	"log"

	"github.com/untrik/FromSkateToZOH/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=untrick  password=4thtgeirf_2001 port=5432 connect_timeout=60 dbname=dbFromSkateToZOH sslmode=disable"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print("Проблема с подключением к базе данных:", err)
		return
	}
	log.Println("Подключение к базе данных установлено.")
	err = DB.AutoMigrate(&models.Admin{}, &models.Student{}, &models.Product{}, &models.Event{}, &models.Reward{}, &models.EventParticipant{})
	if err != nil {
		log.Print("ошибка миграции: ", err)
		return
	}
}
