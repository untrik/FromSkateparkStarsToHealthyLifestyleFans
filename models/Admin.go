package models

type Admin struct {
	AdminId   uint `json:"admin_id" gorm:"primaryKey;autoIncrement"`
	UserID    uint `json:"user_id" gorm:"not null;unique"`
	User      User `json:"user" gorm:"foreignKey:UserID"`
	CreatedBy uint `json:"created_by"`
}
