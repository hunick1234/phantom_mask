package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hunick1234/phantom_mask/application"
	"github.com/hunick1234/phantom_mask/application/query"
	"github.com/hunick1234/phantom_mask/config"
	"github.com/hunick1234/phantom_mask/controller"
	"github.com/hunick1234/phantom_mask/infrastructure/repository"
	"github.com/hunick1234/phantom_mask/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	server := gin.Default()
	server.Use(gin.Logger())
	utils.Init()

	//set db
	cfg := config.LoadConfig()
	dsn := cfg.DB.ToDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//set pharmacy api
	pharmacyQuery := query.NewPharmacyQuery(db)
	ctr := controller.NewPharmacyController(pharmacyQuery)
	controller.SetPharmacyRouter(server, ctr)

	//set user api
	maskRepo := repository.NewMaskRepo(db)
	userRepo := repository.NewUserRepo(db)
	transactionRepo := repository.NewTransactionRepo(db)
	pharmacyRepo := repository.NewPharmacyRepo(db)
	userQuery := query.NewUserQuery(db)
	puchaseSevice := application.NewPurchaseService(
		userRepo,
		pharmacyRepo,
		transactionRepo,
		maskRepo,
	)
	userService := application.NewUserService(puchaseSevice)
	userController := controller.NewUserController(userService, userQuery)
	controller.SetUserRouter(server, userController)

	//set transaction api
	transactionQuery := query.NewTransactionQueryService(db)
	transactionController := controller.NewTransactionController(transactionQuery)
	controller.SetTransactionRouter(server, transactionController)

	//set search api
	maskQuery := query.NewMasksQuery(db)
	searchController := controller.NewSearchController(pharmacyQuery, maskQuery)
	controller.SetSearchRouter(server, searchController)

	err = server.Run(":" + cfg.Port)
	if err != nil {
		panic("failed to start server err: " + err.Error())
	}

}
