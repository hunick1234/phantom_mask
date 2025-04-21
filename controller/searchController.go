package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunick1234/phantom_mask/application/query"
)

type SearchController struct {
	pharmacyQuery *query.PharmacyQueryService
	maskQuery     *query.MasksQuery
}

func NewSearchController(pharmacyQuery *query.PharmacyQueryService, maskQuery *query.MasksQuery) *SearchController {
	return &SearchController{pharmacyQuery, maskQuery}
}

func (s *SearchController) Search(c *gin.Context) {
	var keyword struct {
		Keyword string `form:"q" binding:"required"`
	}
	if err := c.ShouldBindQuery(&keyword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phResults, _ := s.pharmacyQuery.SearchPharmaciesByKeyword(query.PharmacySearchQuery{Keyword: keyword.Keyword})
	maskResults, _ := s.maskQuery.SearchMasksByKeyword(query.SearchMasksQuery{Keyword: keyword.Keyword})

	type SearchResult struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"` // pharmacy 或 mask
	}

	var result []SearchResult
	for _, pharmacy := range phResults {
		result = append(result, SearchResult{
			ID:   pharmacy.ID,
			Name: pharmacy.Name,
			Type: pharmacy.Type,
		})
	}
	for _, mask := range maskResults {
		result = append(result, SearchResult{
			ID:   mask.ID,
			Name: mask.Name,
			Type: mask.Type,
		})
	}
	c.JSON(http.StatusOK, result)
}

func SetSearchRouter(router *gin.Engine, controller *SearchController) {
	router.GET("/api/search", controller.Search)
}
