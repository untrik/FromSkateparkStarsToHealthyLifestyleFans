package models

type Product struct {
	ProductId   uint   `json:"product_id" gorm:"primaryKey"`
	Price       uint   `json:"price"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_URL"`
}
