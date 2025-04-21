package query

import "gorm.io/gorm"

type MasksQuery struct {
	db *gorm.DB
}

type SearchMasksQuery struct {
	Keyword string
}

type SearchMasksDTO struct {
	ID   string
	Type string
	Name string
}

func NewMasksQuery(db *gorm.DB) *MasksQuery {
	return &MasksQuery{
		db: db,
	}
}

func (s *MasksQuery) SearchMasksByKeyword(q SearchMasksQuery) ([]SearchMasksDTO, error) {
	var result []SearchMasksDTO

	sql := `
		SELECT 
			masks.id,
			masks.name,
			'mask' AS type
		FROM masks
		WHERE masks.name ILIKE '%' || $1 || '%'
	`

	err := s.db.Raw(sql, q.Keyword).Scan(&result).Error

	return result, err
}
