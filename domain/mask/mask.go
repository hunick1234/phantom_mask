package mask

import "errors"

type Mask struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Name       string  `gorm:"not null" json:"name"`
	Price      float64 `gorm:"not null" json:"price"`
	Stock      int     `gorm:"not null" json:"stock"`
	PharmacyID uint    `gorm:"index" json:"pharmacy_id"` // 所屬藥局
}


func (m *Mask) CanOffer(quantity int)error {
	if m.Stock < quantity {
		return errors.New("insufficient stock")
	}
	return nil
}