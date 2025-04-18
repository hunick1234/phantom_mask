package query

import (
	"testing"

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

	schema := `
	DROP TABLE IF EXISTS masks;
	CREATE TABLE masks (
		id SERIAL PRIMARY KEY,
		name TEXT,
		price FLOAT,
		stock INT,
		pharmacy_id SERIAL REFERENCES pharmacies(id)
	);
	`

	err = db.Exec(schema).Error
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	testData := `
	INSERT INTO masks (id, name, price,stock, pharmacy_id) VALUES
	('1','Mask A', 10.0, '1', '1'),
	('33','Mask c', 25.0, '1', '1'),
	('22','Mask B', 15.0, '1', '1'),
	('44','Mask d', 5.0, '2', '1'),
	
	('3','nick C', 20.0, '1', '2'),
	('4','yuck D', 25.0, '1', '3'),
	('5','hex E', 30.0, '1', '4');
	`
	err = db.Exec(testData).Error
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	return db
}

func TestSearchMasksByKeyword(t *testing.T) {
	db := setupTestMaskDB(t)
	defer db.Exec("DROP TABLE IF EXISTS masks;")

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
			expected: 4,
		},
		{
			keyword:  "ck",
			expected: 2,
		},
		{
			keyword:  "k",
			expected: 6,
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
