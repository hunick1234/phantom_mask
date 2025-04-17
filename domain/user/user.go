package user

import (
	"github.com/hunick1234/phantom_mask/domain/transaction"
)

type User struct {
	ID           uint                      `gorm:"primaryKey"`
	Name         string                    `gorm:"not null"`
	CashBalance  float64                   `gorm:"not null"`
	Transactions []transaction.Transaction `gorm:"foreignKey:UserID"`
}

type PurchaseInfo struct {
	Product struct {
		Name  string
		Price float64
	}
	Store struct {
		Name string
	}
}
