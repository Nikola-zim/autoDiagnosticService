package middlewares

import (
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	userKey = "user"
)

type Auth struct {
	l logger.Interface
}

func NewAuth(l logger.Interface) *Auth {
	return &Auth{
		l: l,
	}
}

// AuthRequired - middlewares для cookie
func (au *Auth) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(userKey)

		if user == nil {
			// Abort the request with the appropriate error code
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		au.l.Info("user_ID: %v \n", user)

		// Передача запроса
		c.Next()
	}
}
