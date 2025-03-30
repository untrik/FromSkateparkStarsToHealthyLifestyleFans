package models

type Reward struct {
	RewardID  uint `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID   uint `json:"event_id"`
	Place     int  `json:"place" gorm:"check:place > 0"`
	Points    int  `json:"points" gorm:"check:points > 0"`
	StudentID uint `json:"student_id"`
}
