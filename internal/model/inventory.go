package model

type BookInventory struct {
	ISBN            string `gorm:"primaryKey"`
	LibID           string
	Title           string
	Authors         string
	Publisher       string
	Version         string
	TotalCopies     uint
	AvailableCopies uint
}
