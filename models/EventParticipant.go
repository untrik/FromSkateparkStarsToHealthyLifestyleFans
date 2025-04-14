package models

type StatusStudent string

const (
	StatusRegistered StatusStudent = "registered"
	StatusApproved   StatusStudent = "approved"
	StatusCancelled  StatusStudent = "cancelled"
)

type EventParticipant struct {
	EventID   uint          `gorm:"primaryKey" json:"event_id"`
	StudentID uint          `gorm:"primaryKey" json:"student_id"`
	Status    StatusStudent `json:"status" gorm:"default:registered"`
}
