package models

type Users struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string
}
