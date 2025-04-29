package models

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string `json:"username" gorm:"unique;not null;size:200"`
	PasswordHash string `json:"-" gorm:"not null;size:200"`
}
