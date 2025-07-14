package cmd

import (
	"ewallet-transaction/external"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/api"
	"ewallet-transaction/internal/interfaces"
	"ewallet-transaction/internal/repository"
	"ewallet-transaction/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func ServeHTTP() {
	d := dependencyInject()

	r := gin.Default()

	r.GET("/health", d.HealthcheckApi.HealthcheckHandlerHttp)

	transactionV1 := r.Group("/transaction/v1")
	transactionV1.POST("/create", d.ValidateToken, d.TransactionApi.CreateTransaction)

	err := r.Run(":" + helpers.GetEnv("PORT", "8080"))
	if err != nil {
		log.Fatal(err)
	}
}

type Dependency struct {
	HealthcheckApi interfaces.IHealthcheckAPI
	TransactionApi interfaces.ITransactionAPI
	External       interfaces.IExternal
}

func dependencyInject() Dependency {
	healthcheckSvc := &services.Healthcheck{}
	healthcheckAPI := &api.Healthcheck{
		HealthcheckServices: healthcheckSvc,
	}

	transactionRepo := &repository.TransactionRepo{
		DB: helpers.DB,
	}

	transactionSvc := &services.TransactionService{
		TransactionRepo: transactionRepo,
	}
	transactionAPI := &api.TransactionAPI{
		TransactionService: transactionSvc,
	}

	external := &external.External{}

	return Dependency{
		HealthcheckApi: healthcheckAPI,
		TransactionApi: transactionAPI,
		External:       external,
	}
}
