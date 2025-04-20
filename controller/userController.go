package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hunick1234/phantom_mask/application"
	"github.com/hunick1234/phantom_mask/application/query"
	"github.com/hunick1234/phantom_mask/controller/middleware"
	"github.com/hunick1234/phantom_mask/utils"
)

type UserController struct {
	service application.UserSevice
	query   *query.UserQueryService
}

func NewUserController(service application.UserSevice, query *query.UserQueryService) *UserController {
	return &UserController{
		service: service,
		query:   query,
	}
}

func (u *UserController) GetTopTransactionUser(c *gin.Context) {
	var query query.TopUsersTransactionQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := u.query.GetTopUsersByTransactionAmount(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func (u *UserController) Purchase(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var query struct {
		PharmacyID uint `json:"pharmacy_id" binding:"required"`
		MaskID     uint `json:"mask_id" binding:"required"`
		Quantity   int  `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := u.service.Purchase(userID.(uint), query.PharmacyID, query.MaskID, query.Quantity)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (u *UserController) Login(c *gin.Context) {
	//no password login now
	var form struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(400, gin.H{"error": "user_id required"})
		return
	}

	token, err := utils.GenerateJWT(form.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "JWT generation failed"})
		return
	}

	c.JSON(200, gin.H{"jwt_token": token})
}

func SetUserRouter(router *gin.Engine, controller *UserController) {
	user := router.Group("/api/users")
	{
		user.POST("/login", controller.Login)
		user.GET("/top_transactions", controller.GetTopTransactionUser)

		protected := user.Group("/me")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/purchase", controller.Purchase)
		}
	}
}
