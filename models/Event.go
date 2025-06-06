package models

import "time"

type Event struct {
	EventId      uint       `json:"event_id" gorm:"primaryKey;autoIncrement"`
	Title        string     `json:"title"`
	Date         time.Time  `json:"date" gorm:"check:date > NOW()"`
	Location     string     `json:"location"`
	Participants []*Student `json:"participants" gorm:"many2many:event_participants;"`
	Rewards      []Reward   `json:"rewards" gorm:"foreignKey:EventID"`
	Description  string     `json:"description"`
	AdminId      uint       `json:"admin_id" gorm:"not null"`
	Admin        Admin      `json:"admin" gorm:"foreignKey:AdminId"`
}
