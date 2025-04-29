package models

type Product struct {
	ProductId   uint   `json:"product_id" gorm:"primaryKey;autoIncrement"`
	Price       uint   `json:"price" gorm:"not null"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	ImageURL    string `json:"image_URL"`
	Quantity    uint   `json:"quantity" gorm:"not null"`
	IsDeleted   bool   `json:"is_deleted" gorm:"not null;default:false"`
}
