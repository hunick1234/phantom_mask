package query

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"slices"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=localhost user=user password=pass dbname=testdb port=5435 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	schema := `
	DROP TABLE IF EXISTS masks;
	DROP TABLE IF EXISTS pharmacies;
	CREATE TABLE pharmacies (
		id SERIAL PRIMARY KEY,
		name TEXT,
		address TEXT,
		opening_hours JSONB
	);
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
	INSERT INTO pharmacies (id, name, address, opening_hours) VALUES 
	('1', 'Pharmacy A', 'Address A', '{"Mon": ["08:00", "20:00"]}'),
	('2', 'Pharmacy B', 'Address B', '{"Mon": ["10:00", "22:00"]}'),
	('3', 'Pharmacy C', 'Address C', '{"Tue": ["09:00", "18:00"]}'),
	('4', 'Pharmacy D', 'Address D', '{"Wed": ["08:00", "20:00"]}'),
	('5', 'Pharmacy E', 'Address E', '{"Thu": ["10:00", "22:00"]}'),
	('6', 'Pharmacy F', 'Address F', '{"Thu": ["09:00", "18:00"], "Sat": ["10:00", "22:00"], "Sun":["10:00", "12:00"]}');
	
	INSERT INTO masks (id, name, price,stock, pharmacy_id) VALUES
	('1','Mask A', 10.0, '1', '1'),
	('33','Mask c', 25.0, '1', '1'),
	('22','Mask B', 15.0, '1', '1'),
	('44','Mask d', 5.0, '2', '1'),
	
	('3','Mask C', 20.0, '1', '2'),
	('4','Mask D', 25.0, '1', '3'),
	('5','Mask E', 30.0, '1', '4');
	`
	err = db.Exec(testData).Error
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	return db
}

func TestGetOpenPharmaciesOfTime(t *testing.T) {
	db := setupTestDB(t)
	service := &PharmacyQueryService{db: db}

	tests := []struct {
		time     string
		day      string
		expected []string
	}{
		{"09:00", "Mon", []string{"Pharmacy A"}},
		{"10:30", "Mon", []string{"Pharmacy A", "Pharmacy B"}},
		{"09:00", "Tue", []string{"Pharmacy C"}},
		{"09:00", "Wed", []string{"Pharmacy D"}},
		{"11:00", "Thu", []string{"Pharmacy E", "Pharmacy F"}},
		{"09:00", "Fri", []string{}},
		{"10:30", "Sat", []string{"Pharmacy F"}},
		{"11:00", "Sun", []string{"Pharmacy F"}},
		{"08:00", "Sun", []string{}},
	}

	for _, tt := range tests {
		query := OpenPharmacieQuery{Time: tt.time, DayOfWeek: tt.day}
		result, err := service.GetOpenPharmaciesOfTime(query)
		if err != nil {
			t.Errorf("[%s %s] unexpected error: %v", tt.day, tt.time, err)
			continue
		}

		var names []string
		for _, p := range result {
			names = append(names, p.Name)
		}

		for _, expectedName := range tt.expected {
			found := false
			for _, name := range names {
				if name == expectedName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("[%s %s] expected pharmacy '%s' not found", tt.day, tt.time, expectedName)
			}
		}
	}
}

func TestGetMasksByPharmacy(t *testing.T) {
	db := setupTestDB(t)

	service := &PharmacyQueryService{db: db}
	tests := []struct {
		pharmacyID uint
		sortBy     string
		expected   int
	}{
		{1, "price", 4},
		{1, "name", 4},
		{2, "price", 1},
		{3, "price", 1},
		{100, "price", 0},
	}

	for _, v := range tests {
		query := PharmacyMasksQuery{
			PharmacyID: v.pharmacyID,
			SortBy:     v.sortBy,
		}
		masks, err := service.GetMasksByPharmacy(query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(masks) != v.expected {
			t.Errorf("expected %d masks, got %d", v.expected, len(masks))
		}

		// Check masks are sorted correctly
		if v.sortBy == "price" {
			for i := range len(masks) - 1 {
				if masks[i].Price > masks[i+1].Price {
					t.Errorf("masks not sorted by price")
				}
			}
		}
		if v.sortBy == "name" {
			for i := range len(masks) - 1 {
				if masks[i].Name > masks[i+1].Name {
					t.Errorf("masks not sorted by name")
				}
			}
		}

	}

}

func TestGetPharmaciesByMaskCount(t *testing.T) {
	db := setupTestDB(t)
	service := &PharmacyQueryService{db: db}

	tests := []struct {
		name       string
		query      FilterMaskCountQuery
		expectIDs  []string
		expectSize int
	}{
		{
			name: "more than 2 masks",
			query: FilterMaskCountQuery{
				MinPrice:   2.0,
				MaxPrice:   11.0,
				Comparison: "more",
				Count:      2,
			},
			expectIDs:  []string{"1"},
			expectSize: 1,
		},
		{
			name: "less than 2 masks",
			query: FilterMaskCountQuery{
				MinPrice:   2.0,
				MaxPrice:   10.0,
				Comparison: "less",
				Count:      2,
			},
			expectIDs:  []string{},
			expectSize: 0,
		},
		{
			name: "more than 0 masks",
			query: FilterMaskCountQuery{
				MinPrice:   1.0,
				MaxPrice:   30.0,
				Comparison: "more",
				Count:      0,
			},
			expectIDs:  []string{"1","2","3","4"},
			expectSize: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetPharmaciesByMaskCount(tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != tt.expectSize {
				t.Errorf("expected %d pharmacies, got %d", tt.expectSize, len(result))
			}
			for _, pharmacy := range result {
				found := slices.Contains(tt.expectIDs, pharmacy.ID)
				if !found {
					t.Errorf("expected pharmacy ID %s not found", pharmacy.ID)
				}
			}
		})
	}
}
