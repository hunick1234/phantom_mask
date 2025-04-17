package query

import (
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
	DayOfWeek string // Mon...Sun（可選）
}

type OpenPharmacieDTO struct {
	ID        string
	Name      string
	Address   string
	OpenHours map[string][2]string // {"Mon": ["08:00", "20:00"], ...}
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
