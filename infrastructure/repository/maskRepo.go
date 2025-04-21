package repository

import (
	"github.com/hunick1234/phantom_mask/domain/mask"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MaskRepoImpl struct {
	db *gorm.DB
}

func NewMaskRepo(db *gorm.DB) mask.MaskRepo {
	return &MaskRepoImpl{db: db}
}

func (r *MaskRepoImpl) FindByID(pharmacyID, maskID uint) (mask.Mask, error) {
	var m mask.Mask
	if err := r.db.Where("id = ? AND pharmacy_id = ?", maskID, pharmacyID).First(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (r *MaskRepoImpl) FindByIDWithTx(tx *gorm.DB, pharmacyID uint, maskID uint) (mask.Mask, error) {
	var m mask.Mask
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND pharmacy_id = ?", maskID, pharmacyID).First(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (r *MaskRepoImpl) SaveWithTx(tx *gorm.DB, mask *mask.Mask) error {
	return tx.Save(mask).Error
}