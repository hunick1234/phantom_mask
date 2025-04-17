package query

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
func TestGetOpenPharmaciesOfTime(t *testing.T) {
	var testData = `
	INSERT INTO pharmacies (id, name, address, opening_hours) VALUES 
	('1', 'Pharmacy A', 'Address A', '{"Mon": ["08:00", "20:00"]}'),
	('2', 'Pharmacy B', 'Address B', '{"Mon": ["10:00", "22:00"]}'),
	('3', 'Pharmacy C', 'Address C', '{"Tue": ["09:00", "18:00"]}'),
	('4', 'Pharmacy D', 'Address D', '{"Wed": ["08:00", "20:00"]}'),
	('5', 'Pharmacy E', 'Address E', '{"Thu": ["10:00", "22:00"]}'),
	('6', 'Pharmacy F', 'Address F', '{"Thu": ["09:00", "18:00"], "Sat": ["10:00", "22:00"], "Sun":["10:00", "12:00"]}')
	`

	// 測試資料庫
	dsn := "host=localhost user=user password=pass dbname=testdb port=5435 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	// 初始化資料表
	err = db.Exec(`DROP TABLE IF EXISTS pharmacies`).Error
	if err != nil {
		t.Fatalf("failed to drop table: %v", err)
	}

	err = db.Exec(`CREATE TABLE pharmacies (
		id TEXT PRIMARY KEY,
		name TEXT,
		address TEXT,
		opening_hours JSONB
	)`).Error
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	err = db.Exec(testData).Error
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	service := &PharmacyQueryService{db: db}

	// 多組測資
	tests := []struct {
		time      string
		day       string
		expected  []string
	}{
		{"09:00", "Mon", []string{"Pharmacy A"}},
		{"10:30", "Mon", []string{"Pharmacy A", "Pharmacy B"}},
		{"09:00", "Tue", []string{"Pharmacy C"}},
		{"09:00", "Wed", []string{"Pharmacy D"}},
		{"11:00", "Thu", []string{"Pharmacy E","Pharmacy F"}},
		{"09:00", "Fri", []string{}},
		{"10:30", "Sat", []string{"Pharmacy F"}},
		{"11:00", "Sun", []string{"Pharmacy F"}},
		{"08:00", "Sun", []string{}},
	}

	for _, tt := range tests {
		query := OpenPharmacieQuery{
			Time:      tt.time,
			DayOfWeek: tt.day,
		}
		result, err := service.GetOpenPharmaciesOfTime(query)
		if err != nil {
			t.Errorf("[%s %s] unexpected error: %v", tt.day, tt.time, err)
			continue
		}

		if len(result) != len(tt.expected) {
			t.Errorf("[%s %s] expected %d pharmacies, got %d", tt.day, tt.time, len(tt.expected), len(result))
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
				t.Errorf("[%s %s] expected pharmacy '%s' not found in result", tt.day, tt.time, expectedName)
			}
		}
	}
}
