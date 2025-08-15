package models

type Stock struct {
	ID      uint64  `gorm:"primaryKey;autoIncrement"`
	Balance float64 `gorm:"type:numeric(20,4);not null"`
	Reserve float64 `gorm:"type:numeric(20,4);not null"` // ต้องมี column ใน DB
	OnHand  float64 `gorm:"type:numeric(20,4);not null"`
}
