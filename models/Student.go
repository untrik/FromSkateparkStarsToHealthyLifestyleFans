package models

type Student struct {
	StudentId  uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string   `json:"name" gorm:"not null;size:200"`
	LastName   string   `json:"last_name" gorm:"not null;size:200"`
	SecondName string   `json:"second_name" gorm:"size:200"`
	Faculty    string   `json:"faculty" gorm:"not null"`
	Points     uint     `json:"points" gorm:"default:0"`
	UserID     uint     `json:"user_id" gorm:"not null;unique"`
	User       User     `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Rewards    []Reward `json:"rewards" gorm:"foreignKey:StudentID"`
}
