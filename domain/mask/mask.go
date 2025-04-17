package mask

type Mask struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Name       string  `gorm:"not null" json:"name"`
	Price      float64 `gorm:"not null" json:"price"`
	Stock      int     `gorm:"not null" json:"stock"`
	PharmacyID uint    `gorm:"index" json:"pharmacy_id"` // 所屬藥局
}
