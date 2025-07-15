package cmd

import (
	"ewallet-transaction/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Dependency) ValidateToken(c *gin.Context) {
	var (
		log = helpers.Logger
	)
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized empty", nil)
		c.Abort()
		return
	}

	tokenData, err := d.External.ValidateToken(c.Request.Context(), auth)
	if err != nil {
		log.Error(err)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized empty", nil)
		c.Abort()
		return
	}

	tokenData.Token = auth

	c.Set("token", tokenData)

	c.Next()
}
