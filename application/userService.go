package application

type UserSevice interface {
	Purchase(userID, pharmacyID, maskID uint, quanity int) (PurchaseResult, error)
}

type userService struct {
	puchase PurchaseService
}

func NewUserService(purchase PurchaseService) UserSevice {
	return &userService{
		puchase: purchase,
	}
}

func (u *userService) Purchase(userID, pharmacyID, maskID uint, quanity int) (PurchaseResult, error) {
	result, err := u.puchase.Execute(userID, pharmacyID, maskID, quanity)
	if err != nil {
		return PurchaseResult{}, err
	}

	return result, nil
}
