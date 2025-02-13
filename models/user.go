package models

type Users struct {
	ID            string `gorm:"primaryKey"`
	Name          string
	Email         string
	ContactNumber string
	Role          string
	LibID         string
}
