package transaction

import "time"

const (
	StatusPending  = "pending"
	StatusSuccess  = "success"
	StatusFailed   = "failed"
)
type Transaction struct {
	ID                uint              `gorm:"primaryKey"`
	UserID            uint              `gorm:"index"`
	PharmacyID        uint              `gorm:"index"`
	TransactionDate   time.Time         `gorm:"not null"`
	TransactionAmount float64           `gorm:"not null"` // total amount
	Status            string            `gorm:"not null"` // "pending", "succese", "canceled", "failed"
	Message           string            `gorm:"not null"` //  status message
	Items             []TransactionItem `gorm:"foreignKey:TransactionID"`
}

type TransactionItem struct {
	ID            uint    `gorm:"primaryKey"`
	TransactionID uint    `gorm:"index"`
	MaskID        uint    `gorm:"index"`
	Quantity      int     `gorm:"not null"`
	PricePerUnit  float64 `gorm:"not null"`
}
