package pharmacy

import (
	"github.com/hunick1234/phantom_mask/domain/mask"
	"github.com/hunick1234/phantom_mask/utils"
)

type Pharmacy struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	Name         string            `json:"name" gorm:"not null;unique"`
	OpeningHours utils.OpenDayTime `json:"openingHours" gorm:"column:opening_hours;type:jsonb"`
	CashBalance  float64           `json:"cashBalance" gorm:"not null"`
	Masks        []mask.Mask       `json:"masks" gorm:"foreignKey:PharmacyID"`
}

func (p *Pharmacy) AddCash(price float64) {
	p.CashBalance += price
}
