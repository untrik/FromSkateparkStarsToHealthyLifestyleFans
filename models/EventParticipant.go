package models

type EventParticipant struct {
	EventID   uint   `gorm:"primaryKey" json:"event_id"`
	StudentID uint   `gorm:"primaryKey" json:"student_id"`
	Status    string `json:"status"`
	Attended  bool   `gorm:"default:false" json:"attended"`
}
