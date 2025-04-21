package query

import (
	"testing"

	"github.com/hunick1234/phantom_mask/domain/mask"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestMaskDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=localhost user=user password=pass dbname=testdb port=5435 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	err = db.Exec("DROP TABLE IF EXISTS masks CASCADE").Error
	if err != nil {
		t.Fatalf("failed to drop tables: %v", err)
	}
	err = db.AutoMigrate(&mask.Mask{})

	if err != nil {
		t.Fatalf("failed to auto migrate schema: %v", err)
	}
	// Insert test data
	masks := []mask.Mask{
		{Name: "Mask A", Price: 10.0, Stock: 1},
		{Name: "Mask B", Price: 15.0, Stock: 1},
		{Name: "Mask C", Price: 20.0, Stock: 1},
		{Name: "Mask D", Price: 25.0, Stock: 1},
		{Name: "Mask E", Price: 30.0, Stock: 1},
		{Name: "Mask F", Price: 35.0, Stock: 1},
		{Name: "Mak ck", Price: 40.0, Stock: 1},
	}
	err = db.Create(&masks).Error
	if err != nil {
		t.Fatalf("failed to insert masks: %v", err)
	}

	return db
}

func TestSearchMasksByKeyword(t *testing.T) {
	db := setupTestMaskDB(t)
	service := &MasksQuery{db: db}

	tests := []struct {
		keyword  string
		expected int
	}{
		{
			keyword:  "Mask A",
			expected: 1,
		},
		{
			keyword:  "Mask",
			expected: 6,
		},
		{
			keyword:  "ck",
			expected: 1,
		},
		{
			keyword:  "k",
			expected: 7,
		},
	}

	for _, test := range tests {
		t.Run(test.keyword, func(t *testing.T) {
			query := SearchMasksQuery{
				Keyword: test.keyword,
			}
			result, err := service.SearchMasksByKeyword(query)
			if err != nil {
				t.Fatalf("failed to search masks: %v", err)
			}

			if len(result) != test.expected {
				t.Errorf("expected %d results, got %d", test.expected, len(result))
			}
		})
	}
}
