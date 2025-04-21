package query

import (
	"testing"

	"slices"

	"github.com/hunick1234/phantom_mask/domain/mask"
	"github.com/hunick1234/phantom_mask/domain/pharmacy"
	"github.com/hunick1234/phantom_mask/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=localhost user=user password=pass dbname=testdb port=5435 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	err = db.Exec("DROP TABLE IF EXISTS masks, pharmacies CASCADE").Error
	if err != nil {
		t.Fatalf("failed to drop tables: %v", err)
	}
	err = db.AutoMigrate(&pharmacy.Pharmacy{}, &mask.Mask{})
	if err != nil {
		t.Fatalf("failed to auto migrate schema: %v", err)
	}

	// Insert test data
	pharmacies := []pharmacy.Pharmacy{
		{ID: 1, Name: "Pharmacy A", OpeningHours: utils.OpenDayTime{"Mon": {"08:00", "20:00"}}},
		{ID: 2, Name: "Pharmacy B", OpeningHours: utils.OpenDayTime{"Mon": {"10:00", "22:00"}}},
		{ID: 3, Name: "Pharmacy C", OpeningHours: utils.OpenDayTime{"Tue": {"09:00", "18:00"}}},
		{ID: 4, Name: "Pharmacy D", OpeningHours: utils.OpenDayTime{"Wed": {"08:00", "20:00"}}},
		{ID: 5, Name: "Pharmacy E", OpeningHours: utils.OpenDayTime{"Thu": {"10:00", "22:00"}}},
		{ID: 6, Name: "Pharmacy F", OpeningHours: utils.OpenDayTime{"Thu": {"09:00", "18:00"}, "Sat": {"10:00", "22:00"}, "Sun": {"10:00", "12:00"}}},
	}
	err = db.Create(&pharmacies).Error
	if err != nil {
		t.Fatalf("failed to insert pharmacies: %v", err)
	}

	masks := []mask.Mask{

		{Name: "Mask A", Price: 2.0, Stock: 1, PharmacyID: 1},
		{Name: "Mask B", Price: 10.0, Stock: 1, PharmacyID: 1},
		{Name: "Mask C", Price: 10.0, Stock: 1, PharmacyID: 1},
		{Name: "Mask D", Price: 25.0, Stock: 1, PharmacyID: 1},
		{Name: "Mask C", Price: 20.0, Stock: 1, PharmacyID: 2},
		{Name: "Mask D", Price: 25.0, Stock: 1, PharmacyID: 3},
		{Name: "Mask E", Price: 30.0, Stock: 1, PharmacyID: 4},
		{Name: "Mask F", Price: 35.0, Stock: 1, PharmacyID: 5},
		{Name: "Mask G", Price: 40.0, Stock: 1, PharmacyID: 6},
	}
	err = db.Create(&masks).Error
	if err != nil {
		t.Fatalf("failed to insert masks: %v", err)
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
			expectIDs:  []string{"1", "2", "3", "4"},
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

func TestSearchPhamaciesByKeyword(t *testing.T) {
	db := setupTestDB(t)
	service := &PharmacyQueryService{db: db}

	tests := []struct {
		keyword      string
		findExpected int
	}{
		{"Mask A", 0},
		{"ar", 6},
	}

	for _, tt := range tests {
		query := PharmacySearchQuery{Keyword: tt.keyword}
		result, err := service.SearchPharmaciesByKeyword(query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var names []string
		for _, mask := range result {
			names = append(names, mask.Name)
		}
		if len(names) != tt.findExpected {
			t.Errorf("expected %d masks, got %d", tt.findExpected, len(names))
		}
	}
}
