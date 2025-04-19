package pharmacy

import "gorm.io/gorm"

type PharmacyRepo interface {
	Create(pharmacy *Pharmacy) error
	Save(pharmacy *Pharmacy) error
	FindByID(id uint) (Pharmacy, error)

	// Transaction-aware
	SaveWithTx(tx *gorm.DB, pharmacy *Pharmacy) error
	FindByIDWithTx(tx *gorm.DB, id uint) (Pharmacy, error)
}
