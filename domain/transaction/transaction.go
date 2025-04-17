package transaction

import "time"

type Transaction struct {
	ID                uint    `gorm:"primaryKey"`
	UserID            uint    `gorm:"index"`
	PharmacyID        uint    `gorm:"index"`
	MaskID     uint    `gorm:"index"`
	TransactionAmount float64 `gorm:"not null"`
	TransactionDate   time.Time
}
