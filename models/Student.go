package models

type Student struct {
	StudentId  uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string   `json:"name" gorm:"not null"`
	LastName   string   `json:"last_name" gorm:"not null"`
	SecondName string   `json:"second_name"`
	Faculty    string   `json:"faculty" gorm:"not null"`
	Points     uint     `json:"points" gorm:"default:0"`
	User       User     `json:"user" gorm:"foreignKey:UserID"`
	UserID     uint     `json:"user_id" gorm:"not null;unique"`
	Rewards    []Reward `json:"rewards" gorm:"foreignKey:StudentID"`
}
