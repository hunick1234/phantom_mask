package query

import (
	"testing"
	"time"

	"github.com/hunick1234/phantom_mask/domain/transaction"
	"github.com/hunick1234/phantom_mask/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestUserDB(t *testing.T) *gorm.DB {
	t.Helper()
	t.Helper()

	dsn := "host=localhost user=user password=pass dbname=testdb port=5435 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	err = db.Exec("DROP TABLE IF EXISTS transactions, users CASCADE").Error
	if err != nil {
		t.Fatalf("failed to drop tables: %v", err)
	}

	err = db.AutoMigrate(&user.User{}, &transaction.Transaction{})
	if err != nil {
		t.Fatalf("failed to auto migrate schema: %v", err)
	}

	// 建立使用者
	users := []user.User{
		{ID: 1, Name: "Alice", CashBalance: 100.0},
		{ID: 2, Name: "Bob", CashBalance: 200.0},
		{ID: 3, Name: "Charlie", CashBalance: 300.0},
		{ID: 4, Name: "David", CashBalance: 400.0},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatalf("failed to create users: %v", err)
	}

	// 建立交易
	transactions := []transaction.Transaction{
		{ID: 1, UserID: 1, TransactionDate: parseTime("2022-01-01"), TransactionAmount: 50.0, Status: "success"},
		{ID: 2, UserID: 1, TransactionDate: parseTime("2023-01-02"), TransactionAmount: 30.0, Status: "success"},
		{ID: 9, UserID: 1, TransactionDate: parseTime("2025-01-03"), TransactionAmount: 100.0, Status: "success"},
		{ID: 3, UserID: 2, TransactionDate: parseTime("2023-01-01"), TransactionAmount: 20.0, Status: "success"},
		{ID: 4, UserID: 2, TransactionDate: parseTime("2023-01-03"), TransactionAmount: 110.0, Status: "success"},
		{ID: 10, UserID: 2, TransactionDate: parseTime("2025-01-04"), TransactionAmount: 110.0, Status: "success"},
		{ID: 5, UserID: 3, TransactionDate: parseTime("2023-01-02"), TransactionAmount: 100.0, Status: "success"},
		{ID: 6, UserID: 3, TransactionDate: parseTime("2023-01-04"), TransactionAmount: 70.0, Status: "success"},
		{ID: 8, UserID: 4, TransactionDate: parseTime("2025-01-05"), TransactionAmount: 90.0, Status: "success"},
	}

	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to create transactions: %v", err)
	}

	return db
}

func parseTime(date string) time.Time {
	t, _ := time.Parse("2006-01-02", date)
	return t
}

func TestGetTopUsersByTransactionAmount(t *testing.T) {
	db := setupTestUserDB(t)
	service := &UserQueryService{db: db}

	tests := []struct {
		name     string
		query    TopUsersTransactionQuery
		expected []TopUserDTO
	}{
		{
			name: "Test Case 1",
			query: TopUsersTransactionQuery{
				StartDate: parseTime("2023-01-01"),
				EndDate:   parseTime("2023-12-31"),
				Top:       1,
			},
			expected: []TopUserDTO{
				{UserID: "3", Name: "Charlie", TotalAmount: 170.0},
			},
		},
		{
			name: "Test Case 2",
			query: TopUsersTransactionQuery{
				StartDate: parseTime("2023-01-01"),
				EndDate:   parseTime("2023-12-31"),
				Top:       100,
			},
			expected: []TopUserDTO{
				{UserID: "3", Name: "Charlie", TotalAmount: 170.0},
				{UserID: "2", Name: "Bob", TotalAmount: 130.0},
				{UserID: "1", Name: "Alice", TotalAmount: 30.0},
			},
		},
		{
			name: "Test Case 3",
			query: TopUsersTransactionQuery{
				StartDate: parseTime("2023-01-01"),
				EndDate:   parseTime("2023-12-31"),
				Top:       0,
			},
			expected: []TopUserDTO{},
		},
		{
			name: "Test Case 4",
			query: TopUsersTransactionQuery{
				StartDate: parseTime("2025-01-01"),
				EndDate:   parseTime("2025-02-02"),
				Top:       5,
			},
			expected: []TopUserDTO{
				{UserID: "2", Name: "Bob", TotalAmount: 110.0},
				{UserID: "1", Name: "Alice", TotalAmount: 100.0},
				{UserID: "4", Name: "David", TotalAmount: 90.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetTopUsersByTransactionAmount(tt.query)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d results, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if result[i].UserID != expected.UserID || result[i].Name != expected.Name || result[i].TotalAmount != expected.TotalAmount {
					t.Errorf("expected %v, got %v", expected, result[i])
				}
			}
		})
	}

}
