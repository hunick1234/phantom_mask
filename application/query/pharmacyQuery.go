package query

import (
	"github.com/hunick1234/phantom_mask/utils"
	"gorm.io/gorm"
)

type PharmacyQueryService struct {
	db *gorm.DB
}

func NewPharmacyQuery() *PharmacyQueryService {
	return &PharmacyQueryService{}
}

type OpenPharmacieQuery struct {
	Time      string // 格式 HH:MM
	DayOfWeek string
}

type OpenPharmacieDTO struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Address   string            `json:"address"`
	OpenHours utils.OpenDayTime `json:"opening_hours" gorm:"column:opening_hours;type:jsonb"`
}

type PharmacyMasksQuery struct {
	PharmacyID uint
	SortBy     string
}
type PharmacyMasksDTO struct {
	Id    string
	Name  string
	Price float64
	Stock int
}

func (p *PharmacyQueryService) GetOpenPharmaciesOfTime(q OpenPharmacieQuery) ([]OpenPharmacieDTO, error) {
	var result []OpenPharmacieDTO

	// TODO:check query is valid
	day := q.DayOfWeek
	timeStr := q.Time

	sql := `
	SELECT *
	FROM pharmacies
	WHERE (
  		$2 >= opening_hours-> $1->>0 AND
  		$2 <= opening_hours-> $1->>1
	);`

	if err := p.db.Raw(sql, day, timeStr).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PharmacyQueryService) GetMasksByPharmacy(query PharmacyMasksQuery) ([]PharmacyMasksDTO, error) {
	var result []PharmacyMasksDTO

	// TODO:check query is valid
	sql := `
	SELECT id, name, price, stock
	FROM masks
	WHERE pharmacy_id = $1
	ORDER BY ` + query.SortBy + ` ASC;`

	if err := p.db.Raw(sql, query.PharmacyID).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
