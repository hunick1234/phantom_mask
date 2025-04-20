package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hunick1234/phantom_mask/application/query"
)

type PharmacyController struct {
	query *query.PharmacyQueryService
}

func NewPharmacyController(q *query.PharmacyQueryService) *PharmacyController {
	return &PharmacyController{
		query: q,
	}
}

func (p *PharmacyController) GetOpenPharmacies(c *gin.Context) {
	var query query.OpenPharmacieQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := p.query.GetOpenPharmaciesOfTime(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (p *PharmacyController) GetMasksByPharmacy(c *gin.Context) {
	var query query.PharmacyMasksQuery
	
	idParam := c.Param("pharmacy_id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pharmacy_id"})
		return
	}
	query.PharmacyID = uint(id)

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := p.query.GetMasksByPharmacy(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (p *PharmacyController) FilterPharmaciesByMaskCount(c *gin.Context) {
	var query query.FilterMaskCountQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := p.query.GetPharmaciesByMaskCount(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func SetPharmacyRouter(r *gin.Engine, ctr *PharmacyController) {
	api := r.Group("/api/pharmacies")
	api.GET("/open", ctr.GetOpenPharmacies)
	api.GET("/:pharmacy_id/masks", ctr.GetMasksByPharmacy)
	api.GET("/filter_by_mask_count", ctr.FilterPharmaciesByMaskCount)
}
