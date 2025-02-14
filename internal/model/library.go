package model

type Library struct {
	ID   string `gorm:"primaryKey"`
	Name string
}
