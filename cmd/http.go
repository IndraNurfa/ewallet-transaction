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
	transactionV1.POST("/refund", d.ValidateToken, d.TransactionApi.RefundTransaction)
	transactionV1.PUT("/update-status/:reference", d.ValidateToken, d.TransactionApi.UpdateStatusTransaction)
	transactionV1.GET("/:reference", d.ValidateToken, d.TransactionApi.GetTransactionDetail)
	transactionV1.GET("/", d.ValidateToken, d.TransactionApi.GetTransaction)

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

	external := &external.External{}

	transactionRepo := &repository.TransactionRepo{
		DB: helpers.DB,
	}

	transactionSvc := &services.TransactionService{
		TransactionRepo: transactionRepo,
		External:        external,
	}
	transactionAPI := &api.TransactionAPI{
		TransactionService: transactionSvc,
	}

	return Dependency{
		HealthcheckApi: healthcheckAPI,
		TransactionApi: transactionAPI,
		External:       external,
	}
}
