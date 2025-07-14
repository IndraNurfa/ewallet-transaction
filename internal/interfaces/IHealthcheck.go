package interfaces

import "github.com/gin-gonic/gin"

type IHealthcheckAPI interface {
	HealthcheckHandlerHttp(c *gin.Context)
}

type IHealthcheckServices interface {
	HealthcheckServices() (string, error)
}

type IHealthcheckRepo interface {
}
