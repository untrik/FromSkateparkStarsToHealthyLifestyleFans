package models

type Order struct {
	OrderId   uint    `json:"order_id" gorm:"primaryKey;autoIncrement"`
	StudentId uint    `json:"student_id" gorm:"not null"`
	ProductId uint    `json:"product_id" gorm:"not null"`
	Quantity  uint    `json:"quantity"`
	Student   Student `gorm:"foreignKey:StudentId;references:StudentId"`
	Product   Product `gorm:"foreignKey:ProductId;references:ProductId"`
}
