package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthenticateAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"details": "Fix me for login"})
}
