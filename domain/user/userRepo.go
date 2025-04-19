package user

import "gorm.io/gorm"

type UserRepo interface {
	Create(user *User) error
	Save(user *User) error
	FindByID(id uint) (User, error)

	// Transaction-aware 方法
	FindByIDWithTx(tx *gorm.DB, id uint) (User, error)
	SaveWithTx(tx *gorm.DB, user *User) error
}
