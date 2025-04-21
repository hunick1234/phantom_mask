package query

import (
	"github.com/hunick1234/phantom_mask/utils"
	"gorm.io/gorm"
)

type PharmacyQueryService struct {
	db *gorm.DB
}

func NewPharmacyQuery(db *gorm.DB) *PharmacyQueryService {
	return &PharmacyQueryService{
		db: db,
	}
}

type OpenPharmacieQuery struct {
	Time      string `form:"time" binding:"required,timeformat"`
	DayOfWeek string `form:"day_of_week" binding:"required,dayofweek"`
}

type OpenPharmacieDTO struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Address   string            `json:"address"`
	OpenHours utils.OpenDayTime `json:"opening_hours" gorm:"column:opening_hours;type:jsonb"`
}

type PharmacyMasksQuery struct {
	PharmacyID uint   
	SortBy     string `form:"sort_by" binding:"required,sortby"`
}

type PharmacyMasksDTO struct {
	Id    string
	Name  string
	Price float64
	Stock int
}

type FilterMaskCountQuery struct {
	MinPrice   float64 `form:"min_price" binding:"required,min=0"`
	MaxPrice   float64 `form:"max_price" binding:"required,min=0,gtfield=MinPrice"`
	Comparison string  `form:"comparison" binding:"required,comparison"`
	Count      int     `form:"count" binding:"required,min=0"`
}

type PharmacyMaskCountDTO struct {
	ID        string
	Name      string
	MaskCount int
}

type PharmacySearchQuery struct {
	Keyword string
}

type PharmacySearchDTO struct {
	ID   string
	Type string
	Name string
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

func (p *PharmacyQueryService) GetPharmaciesByMaskCount(query FilterMaskCountQuery) ([]PharmacyMaskCountDTO, error) {
	var result []PharmacyMaskCountDTO

	comparisonOp := ">"
	if query.Comparison == "less" {
		comparisonOp = "<"
	}
	sql := `
		SELECT 
			pharmacies.id AS id,
			pharmacies.name AS name,
			SUM(masks.stock) AS mask_count
		FROM pharmacies
		JOIN masks ON masks.pharmacy_id = pharmacies.id
		WHERE masks.price BETWEEN $1 AND $2
		GROUP BY pharmacies.id, pharmacies.name
		HAVING SUM(masks.stock) ` + comparisonOp + ` $3
	`

	if err := p.db.Raw(sql, query.MinPrice, query.MaxPrice, query.Count).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PharmacyQueryService) SearchPharmaciesByKeyword(query PharmacySearchQuery) ([]PharmacySearchDTO, error) {
	var result []PharmacySearchDTO

	sql := `
	SELECT 
		pharmacies.id,
		pharmacies.name,
		'pharmacy' AS type
	FROM pharmacies
	WHERE name ILIKE '%' || $1 || '%'
	`

	if err := p.db.Raw(sql, query.Keyword).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
