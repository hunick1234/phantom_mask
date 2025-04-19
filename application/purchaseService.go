package application

import (
	"context"
	"fmt"
	"time"

	"github.com/hunick1234/phantom_mask/domain/mask"
	"github.com/hunick1234/phantom_mask/domain/pharmacy"
	"github.com/hunick1234/phantom_mask/domain/transaction"
	"github.com/hunick1234/phantom_mask/domain/user"
	"gorm.io/gorm"
)

type PurchaseService interface {
	Execute(userID, pharmacyID, maskID uint, quantity int) (PurchaseResult, error)
}

type purchaseService struct {
	userRepo        user.UserRepo
	storeRepo       pharmacy.PharmacyRepo
	transactionRepo transaction.TransactionRepo
	maskRepo        mask.MaskRepo
}

type PurchaseResult struct {
	TranslactionID uint
	Status         string
}

func NewPurchaseService(userRepo user.UserRepo, storeRepo pharmacy.PharmacyRepo, transactionRepo transaction.TransactionRepo, maskRepo mask.MaskRepo) PurchaseService {
	return &purchaseService{
		userRepo:        userRepo,
		storeRepo:       storeRepo,
		transactionRepo: transactionRepo,
		maskRepo:        maskRepo,
	}
}

func (p *purchaseService) Execute(userID, storeID, productID uint, quantity int) (PurchaseResult, error) {
	var result PurchaseResult

	// Step 1: 開始交易，建立交易紀錄
	t := &transaction.Transaction{
		UserID:            userID,
		PharmacyID:        storeID,
		TransactionAmount: 0,
		TransactionDate:   time.Now(),
		Status:            transaction.StatusPending,
		Items: []transaction.TransactionItem{
			{
				MaskID:       productID,
				Quantity:     quantity,
				PricePerUnit: 0,
			},
		},
	}
	if err := p.transactionRepo.Create(t); err != nil {
		result.Status = "failed"
		return result, err
	}
	result.TranslactionID = t.ID

	// Step 2: 預檢，查詢使用者、藥局、商品
	user, err := p.userRepo.FindByID(userID)
	if err != nil {
		result.Status = "failed"
		return result, p.failTransaction(t, err)
	}

	if _, err = p.storeRepo.FindByID(storeID); err != nil {
		result.Status = "failed"
		return result, p.failTransaction(t, err)
	}

	product, err := p.maskRepo.FindByID(storeID, productID)
	if err != nil {
		result.Status = "failed"
		return result, p.failTransaction(t, err)
	}

	if err := product.CanOffer(quantity); err != nil {
		result.Status = "failed"
		return result, p.failTransaction(t, err)
	}

	totalAmount := product.Price * float64(quantity)
	if err := user.CanAfford(totalAmount); err != nil {
		result.Status = "failed"
		return result, p.failTransaction(t, err)
	}

	// Step 3: 進入交易，扣款、減庫存、更新狀態

	err = p.transactionRepo.WithTx(context.TODO(), func(tx *gorm.DB) error {

		user, err := p.userRepo.FindByIDWithTx(tx, userID)
		if err != nil {
			return err
		}
		store, err := p.storeRepo.FindByIDWithTx(tx, storeID)
		if err != nil {
			return err
		}
		product, err := p.maskRepo.FindByIDWithTx(tx, storeID, productID)
		if err != nil {
			return err
		}

		// Re-check again
		if err := product.CanOffer(quantity); err != nil {
			return err
		}
		if err := user.CanAfford(totalAmount); err != nil {
			return err
		}

		// Perform state mutation
		user.Pay(totalAmount)
		store.AddCash(totalAmount)
		product.Stock -= quantity

		p.userRepo.SaveWithTx(tx, &user)
		p.storeRepo.SaveWithTx(tx, &store)
		p.maskRepo.SaveWithTx(tx, &product)

		// Update transaction status to success
		if err := p.succesesTransaction(t); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		result.Status = "failed"
		return result, p.failTransaction(t, err)
	}

	result.Status = "success"
	return result, nil
}

func (p *purchaseService) succesesTransaction(t *transaction.Transaction) error {
	t.Status = transaction.StatusSuccess
	return p.transactionRepo.Save(t) // 簡化
}

func (p *purchaseService) failTransaction(t *transaction.Transaction, cause error) error {
	if cause != nil {
		t.Status = transaction.StatusFailed
		t.Message = cause.Error()
	}
	saveErr := p.transactionRepo.Save(t)
	if saveErr != nil {
		// 儲存也失敗，就把兩個錯誤合併回傳
		return fmt.Errorf("transaction failed: %v; also failed to save status: %v", cause, saveErr)
	}
	return cause
}
