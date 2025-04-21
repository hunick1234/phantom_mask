package user

import (
	"errors"

	"github.com/hunick1234/phantom_mask/domain/transaction"
)

type User struct {
	ID           uint                      `gorm:"primaryKey"`
	Name         string                    `gorm:"not null"`
	CashBalance  float64                   `gorm:"not null"`
	Transactions []transaction.Transaction `gorm:"foreignKey:UserID"`
}

func (u *User) CanAfford(price float64) error {
	if u.CashBalance < price {
		return errors.New("insufficient funds")
	}
	return nil
}

func (u *User) Pay(price float64) {
	u.CashBalance -= price
}
