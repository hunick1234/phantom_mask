package application

import (
	"fmt"
	"sync"
	"testing"

	"github.com/hunick1234/phantom_mask/domain/mask"
	"github.com/hunick1234/phantom_mask/domain/pharmacy"
	"github.com/hunick1234/phantom_mask/domain/transaction"
	"github.com/hunick1234/phantom_mask/domain/user"
	"github.com/hunick1234/phantom_mask/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupPurchaseTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "host=localhost user=user password=pass dbname=testdb port=5435 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	db.Exec("DROP TABLE IF EXISTS transaction_items, transactions, masks, users, pharmacies CASCADE")
	db.AutoMigrate(&user.User{}, &pharmacy.Pharmacy{}, &mask.Mask{}, &transaction.Transaction{}, &transaction.TransactionItem{})

	return db
}
func TestPurchaseServiceExecute(t *testing.T) {
	db := setupPurchaseTestDB(t)

	cases := []struct {
		name                  string
		userBalance           float64
		pharmacyBalance       float64
		maskStock             int
		maskPrice             float64
		quantity              int
		expectStatus          string
		expectUserBalance     float64
		expectPharmacyBalance float64
		expectMaskStock       int
	}{
		{
			name:        "成功購買",
			userBalance: 100.0, pharmacyBalance: 50.0, maskStock: 10, maskPrice: 10.0, quantity: 2,
			expectStatus: "success", expectUserBalance: 80.0, expectPharmacyBalance: 70.0, expectMaskStock: 8,
		},
		{
			name:        "使用者餘額不足",
			userBalance: 10.0, pharmacyBalance: 50.0, maskStock: 10, maskPrice: 10.0, quantity: 2,
			expectStatus: "failed", expectUserBalance: 10.0, expectPharmacyBalance: 50.0, expectMaskStock: 10,
		},
		{
			name:        "口罩庫存不足",
			userBalance: 100.0, pharmacyBalance: 50.0, maskStock: 1, maskPrice: 10.0, quantity: 2,
			expectStatus: "failed", expectUserBalance: 100.0, expectPharmacyBalance: 50.0, expectMaskStock: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db.Exec("TRUNCATE TABLE transaction_items, transactions, masks, users, pharmacies CASCADE")
			u := user.User{ID: 1, CashBalance: c.userBalance}
			p := pharmacy.Pharmacy{ID: 1, CashBalance: c.pharmacyBalance}
			m := mask.Mask{ID: 1, Name: "Mask A", Stock: c.maskStock, Price: c.maskPrice, PharmacyID: 1}
			db.Create(&u)
			db.Create(&p)
			db.Create(&m)

			ps := NewPurchaseService(
				repository.NewUserRepo(db),
				repository.NewPharmacyRepo(db),
				repository.NewTransactionRepo(db),
				repository.NewMaskRepo(db),
			)

			result, _ := ps.Execute(1, 1, 1, c.quantity)

			var updatedUser user.User
			var updatedPharmacy pharmacy.Pharmacy
			var updatedMask mask.Mask
			var transaction transaction.Transaction
			db.First(&updatedUser, 1)
			db.First(&updatedPharmacy, 1)
			db.First(&updatedMask, 1)
			db.First(&transaction, result.TranslactionID)

			assert.Equal(t, c.expectStatus, result.Status)
			assert.InEpsilon(t, c.expectUserBalance, updatedUser.CashBalance, 0.01)
			assert.InEpsilon(t, c.expectPharmacyBalance, updatedPharmacy.CashBalance, 0.01)
			assert.Equal(t, c.expectMaskStock, updatedMask.Stock)
		})
	}
}

func TestConcurrentPurchase(t *testing.T) {
	db := setupPurchaseTestDB(t)
	db.Exec("TRUNCATE TABLE transaction_items, transactions, masks, users, pharmacies CASCADE")
	ps := NewPurchaseService(
		repository.NewUserRepo(db),
		repository.NewPharmacyRepo(db),
		repository.NewTransactionRepo(db),
		repository.NewMaskRepo(db),
	)

	db.Create(&user.User{ID: 1, CashBalance: 100})
	db.Create(&pharmacy.Pharmacy{ID: 1, CashBalance: 0})
	db.Create(&mask.Mask{ID: 1, Name: "Mask A", Stock: 1, Price: 10.0, PharmacyID: 1}) // 只有 1 片口罩

	var wg sync.WaitGroup
	results := make(chan string, 2)

	wg.Add(2)
	go purchase(ps, 1, &wg, results)
	go purchase(ps, 2, &wg, results)
	wg.Wait()
	close(results)

	fmt.Println("=== 併發購買結果 ===")
	for r := range results {
		fmt.Println(r)
	}
}

func purchase(service PurchaseService, threadID int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	result, err := service.Execute(1, 1, 1, 1)
	if err != nil {
		results <- fmt.Sprintf("Thread %d: Error: %v", threadID, err)
		return
	}
	results <- fmt.Sprintf("Thread %d: Status: %s, TransactionID: %d", threadID, result.Status, result.TranslactionID)
}
