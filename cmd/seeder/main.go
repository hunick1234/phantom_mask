package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hunick1234/phantom_mask/domain/mask"
	"github.com/hunick1234/phantom_mask/domain/pharmacy"
	"github.com/hunick1234/phantom_mask/domain/transaction"
	"github.com/hunick1234/phantom_mask/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RawPharmacy struct {
	Name         string  `json:"name"`
	CashBalance  float64 `json:"cashBalance"`
	OpeningHours string  `json:"openingHours"`
	Masks        []struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	} `json:"masks"`
}

type RawUser struct {
	Name              string  `json:"name"`
	CashBalance       float64 `json:"cashBalance"`
	PurchaseHistories []struct {
		PharmacyName      string  `json:"pharmacyName"`
		MaskName          string  `json:"maskName"`
		TransactionAmount float64 `json:"transactionAmount"`
		TransactionDate   string  `json:"transactionDate"`
	} `json:"purchaseHistories"`
}

func main() {
	var filePath, seedType string
	flag.StringVar(&seedType, "t", "", "Type of data to seed (e.g., 'user' or 'pharmacy')")
	flag.StringVar(&filePath, "p", "", "Path to the JSON file to seed data from")
	flag.Parse()

	if seedType == "" || filePath == "" {
		fmt.Println("Usage: go run main.go -t <type> -p <path>")
		os.Exit(1)
	}

	// Connect to PostgreSQL
	dsn := "host=localhost user=postgres password=postgres dbname=phantom_mask port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Auto migrate schema
	db.AutoMigrate(
		&pharmacy.Pharmacy{},
		&mask.Mask{},
		&user.User{},
		&transaction.Transaction{},
	)

	// Call seeder based on flag
	switch seedType {
	case "pharmacy":
		SeedPharmacies(db, filePath)
	case "user":
		SeedUsers(db, filePath)
	default:
		fmt.Println("Unknown seed type. Use 'pharmacy' or 'user'")
		os.Exit(1)
	}
}

func SeedPharmacies(db *gorm.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	defer file.Close()

	var rawPharmacies []RawPharmacy
	if err := json.NewDecoder(file).Decode(&rawPharmacies); err != nil {
		log.Fatalf("json decode error: %v", err)
	}

	for _, rp := range rawPharmacies {
		pharmacy := pharmacy.Pharmacy{
			Name:         rp.Name,
			CashBalance:  rp.CashBalance,
			OpeningHours: pharmacy.FormateOpeningHours(rp.OpeningHours),
		}

		for _, m := range rp.Masks {
			pharmacy.Masks = append(pharmacy.Masks, mask.Mask{
				Name:  m.Name,
				Price: m.Price,
				Stock: 1, // 預設庫存
			})
		}

		if err := db.Create(&pharmacy).Error; err != nil {
			log.Printf("create pharmacy error: %v", err)
		}
	}
}

func SeedUsers(db *gorm.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	defer file.Close()

	var rawUsers []RawUser
	if err := json.NewDecoder(file).Decode(&rawUsers); err != nil {
		log.Fatalf("decode error: %v", err)
	}

	for _, ru := range rawUsers {
		user := user.User{
			Name:        ru.Name,
			CashBalance: ru.CashBalance,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Printf("create user error: %v", err)
		}

		for _, p := range ru.PurchaseHistories {
			var pharmacy pharmacy.Pharmacy
			var mask mask.Mask

			db.Where("name = ?", p.PharmacyName).First(&pharmacy)
			db.Where("name = ? AND pharmacy_id = ?", p.MaskName, pharmacy.ID).First(&mask)

			tDate, _ := time.Parse("2006-01-02 15:04:05", p.TransactionDate)

			tx := transaction.Transaction{
				UserID:            user.ID,
				PharmacyID:        pharmacy.ID,
				MaskID:            mask.ID,
				TransactionAmount: p.TransactionAmount,
				TransactionDate:   tDate,
			}

			db.Create(&tx)

		}
	}
}
