package mask

import "gorm.io/gorm"

type MaskRepo interface {
	FindByID(pharmacyID, maskID uint) (Mask, error)

	SaveWithTx(tx *gorm.DB, mask *Mask) error
	FindByIDWithTx(tx *gorm.DB,pharmacyID uint, maskID uint) (Mask, error)
}
