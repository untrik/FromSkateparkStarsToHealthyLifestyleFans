package models

type Admin struct {
	AdminId      uint   `json:"admin_id" gorm:"primaryKey;autoIncrement"`
	Username     string `json:"username" gorm:"unique;not null"`
	PasswordHash string `json:"-" gorm:"not null;size:100"`
}
