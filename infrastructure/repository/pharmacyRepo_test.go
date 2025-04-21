package repository

import (
	"testing"

	"github.com/hunick1234/phantom_mask/domain/mask"
	"github.com/hunick1234/phantom_mask/domain/pharmacy"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFindPharmacyByID(t *testing.T) {
	// 初始化測試資料庫（SQLite in-memory）
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to SQLite: %v", err)
	}

	// 建立資料表
	err = db.AutoMigrate(&pharmacy.Pharmacy{}, &mask.Mask{})
	if err != nil {
		t.Fatalf("failed to connect to SQLite: %v", err)
	}

	// 建立測試資料
	p := pharmacy.Pharmacy{
		Name:        "Test Pharmacy",
		CashBalance: 100.0,
		Masks: []mask.Mask{
			{Name: "Mask A", Price: 10.0, Stock: 50},
			{Name: "Mask B", Price: 15.0, Stock: 30},
		},
	}
	err = db.Create(&p).Error
	if err != nil {
		t.Fatalf("failed to connect to SQLite: %v", err)
	}

	// 初始化 repository
	repo := NewPharmacyRepo(db)

	// 呼叫測試目標函數
	result, err := repo.FindByID(p.ID)
	if err != nil {
		t.Fatalf("failed to find pharmacy: %v", err)
	}
	if result.ID != p.ID {
		t.Errorf("expected pharmacy ID %d, got %d", p.ID, result.ID)
	}
}
