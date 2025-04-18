package query

import (
	"testing"
	"time"

	"github.com/hunick1234/phantom_mask/domain/transaction"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupSummaryTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "host=localhost user=user password=pass dbname=testdb port=5435 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	err = db.Exec("DROP TABLE IF EXISTS transaction_items, transactions CASCADE").Error
	if err != nil {
		t.Fatalf("failed to drop tables: %v", err)
	}

	err = db.AutoMigrate(&transaction.Transaction{}, &transaction.TransactionItem{})
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	transactions := []transaction.Transaction{
		{ID: 1, UserID: 1, TransactionDate: parseTime("2025-01-01"), TransactionAmount: 100.0},
		{ID: 2, UserID: 2, TransactionDate: parseTime("2025-01-02"), TransactionAmount: 60.0},
		{ID: 3, UserID: 1, TransactionDate: parseTime("2025-01-03"), TransactionAmount: 50.0},
		{ID: 4, UserID: 3, TransactionDate: parseTime("2025-01-04"), TransactionAmount: 120.0},
		{ID: 5, UserID: 2, TransactionDate: parseTime("2025-01-05"), TransactionAmount: 230.0},
	}
	db.Create(&transactions)

	items := []transaction.TransactionItem{
		{TransactionID: 1, MaskID: 1, Quantity: 2, PricePerUnit: 30.0},
		{TransactionID: 1, MaskID: 2, Quantity: 1, PricePerUnit: 40.0},
		{TransactionID: 2, MaskID: 3, Quantity: 3, PricePerUnit: 20.0},
		{TransactionID: 3, MaskID: 1, Quantity: 1, PricePerUnit: 50.0},
		{TransactionID: 4, MaskID: 2, Quantity: 2, PricePerUnit: 60.0},
		{TransactionID: 5, MaskID: 3, Quantity: 1, PricePerUnit: 70.0},
		{TransactionID: 5, MaskID: 1, Quantity: 2, PricePerUnit: 80.0},
	}
	db.Create(&items)

	return db
}

func TestGetTransactionSummary(t *testing.T) {
	db := setupSummaryTestDB(t)
	service := UserQueryService{DB: db}

	tests := []struct {
		testName  string
		startDate time.Time
		endDate   time.Time
		expected  TransactionSummaryDTO
	}{
		{
			testName:  "test case 1",
			startDate: parseTime("2025-01-01"),
			endDate:   parseTime("2025-01-05"),
			expected: TransactionSummaryDTO{
				TotalMasks:  12,
				TotalAmount: 560,
			},
		},
		{
			testName:  "test case 2",
			startDate: parseTime("2025-01-02"),
			endDate:   parseTime("2025-01-03"),
			expected: TransactionSummaryDTO{
				TotalMasks:  4,
				TotalAmount: 110,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			query := TransactionSummaryQuery{
				StartDate: test.startDate,
				EndDate:   test.endDate,
			}

			result, err := service.GetTransactionSummary(query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.TotalMasks != test.expected.TotalMasks {
				t.Errorf("expected %d masks, got %d", test.expected.TotalMasks, result.TotalMasks)
			}
			if result.TotalAmount != test.expected.TotalAmount {
				t.Errorf("expected %f amount, got %f", test.expected.TotalAmount, result.TotalAmount)
			}
		})
	}
}
