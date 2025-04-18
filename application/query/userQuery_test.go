package query

import (
	"testing"
	"time"

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

	schema := `
	DROP TABLE IF EXISTS transactions;
	DROP TABLE IF EXISTS users;
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		cash_balance FLOAT NOT NULL
	);
	CREATE TABLE transactions (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		transaction_date DATE NOT NULL,
		transaction_amount FLOAT NOT NULL
	);
	`

	err = db.Exec(schema).Error
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	testData := `
	INSERT INTO users (id, name, cash_balance) VALUES
	(1, 'Alice', 100.0),
	(2, 'Bob', 200.0),
	(3, 'Charlie', 300.0),
	(4, 'David', 400.0);
	INSERT INTO TRANSACTIONS (id, user_id, transaction_date, transaction_amount) VALUES
	(1, 1, '2022-01-01', 50.0),
	(2, 1, '2023-01-02', 30.0),
	(9, 1, '2025-01-03', 100.0),
	(3, 2, '2023-01-01', 20.0),
	(4, 2, '2023-01-03', 110.0),
	(10, 2, '2025-01-04', 110.0),
	(5, 3, '2023-01-02', 100.0),
	(6, 3, '2023-01-04', 70.0),
	(8, 4, '2025-01-05', 90.0);
	`
	err = db.Exec(testData).Error
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	return db
}

func parseTime(date string) time.Time {
	t, _ := time.Parse("2006-01-02", date)
	return t
}

func TestGetTopUsersByTransactionAmount(t *testing.T) {
	db := setupTestUserDB(t)
	service := &UserQueryService{DB: db}

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
