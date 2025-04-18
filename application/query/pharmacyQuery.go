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

func (p *PharmacyQueryService) GetOpenPharmaciesOfTime(q OpenPharmacieQuery) ([]OpenPharmacieDTO, error) {
	var result []OpenPharmacieDTO

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

func (p *PharmacyQueryService) GetPharmacySellingMasks() ([]OpenPharmacieDTO, error) {
	var result []OpenPharmacieDTO

	sql := `
	SELECT *
	FROM pharmacies
	WHERE selling_masks = true;`

	if err := p.db.Raw(sql).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
