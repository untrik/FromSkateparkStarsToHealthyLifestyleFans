package models

type Student struct {
	StudentId    uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string   `json:"name" gorm:"not null"`
	LastName     string   `json:"last_name" gorm:"not null"`
	SecondName   string   `json:"second_name"`
	Faculty      string   `json:"faculty" gorm:"not null"`
	Points       uint     `json:"points" gorm:"default:0"`
	Username     string   `json:"username" gorm:"unique;not null"`
	PasswordHash string   `json:"-" gorm:"not null;size:100"`
	Rewards      []Reward `json:"rewards" gorm:"foreignKey:StudentID"`
}
