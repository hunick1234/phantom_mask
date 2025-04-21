package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hunick1234/phantom_mask/application/query"
)

type TransactionController struct {
	query *query.TransactionQueryService
}

func NewTransactionController(query *query.TransactionQueryService) *TransactionController {
	return &TransactionController{
		query: query,
	}
}

func (tc *TransactionController) Summary(c *gin.Context) {
	// Implement the logic to handle the transaction summary
	var query query.TransactionSummaryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	result, err := tc.query.GetTransactionSummary(query)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, result)
}

func SetTransactionRouter(router *gin.Engine, controller *TransactionController) {
	transaction := router.Group("/api/transactions")
	transaction.GET("/summary", controller.Summary)
}
