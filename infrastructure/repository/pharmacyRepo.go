package repository

import (
	"github.com/hunick1234/phantom_mask/domain/pharmacy"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pharmacyRepoImpl struct {
	db *gorm.DB
}

func NewPharmacyRepo(db *gorm.DB) pharmacy.PharmacyRepo {
	return &pharmacyRepoImpl{db: db}
}

func (r *pharmacyRepoImpl) Save(p *pharmacy.Pharmacy) error {
	return r.db.Save(p).Error
}

func (r *pharmacyRepoImpl) Create(p *pharmacy.Pharmacy) error {
	return r.db.Create(p).Error
}

func (r *pharmacyRepoImpl) FindByID(id uint) (pharmacy.Pharmacy, error) {
	var p pharmacy.Pharmacy
	if err := r.db.Preload("Masks").First(&p, id).Error; err != nil {
		return p, err
	}
	return p, nil
}

func (r *pharmacyRepoImpl) SaveWithTx(tx *gorm.DB, p *pharmacy.Pharmacy) error {
	return tx.Save(p).Error
}

func (r *pharmacyRepoImpl) FindByIDWithTx(tx *gorm.DB, id uint) (pharmacy.Pharmacy, error) {
	var p pharmacy.Pharmacy
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&p, id).Error; err != nil {
		return p, err
	}
	return p, nil
}
