package repository

import (
	"github.com/hunick1234/phantom_mask/domain/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) user.UserRepo {
	return &userRepoImpl{db: db}
}

func (r *userRepoImpl) Create(user *user.User) error {
	return r.db.Create(user).Error
}

func (r *userRepoImpl) Save(user *user.User) error {
	return r.db.Save(user).Error
}

func (r *userRepoImpl) FindByID(id uint) (user.User, error) {
	var u user.User
	if err := r.db.First(&u, id).Error; err != nil {
		return u, err
	}
	return u, nil
}

func (r *userRepoImpl) SaveWithTx(tx *gorm.DB, user *user.User) error {
	return tx.Save(user).Error
}

func (r *userRepoImpl) FindByIDWithTx(tx *gorm.DB, id uint) (user.User, error) {
	var u user.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&u, id).Error; err != nil {
		return u, err
	}
	return u, nil
}
