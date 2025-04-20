package middleware

import (
	"gocms/utils"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type JwtMiddleware struct {
	config *viper.Viper
	e      *casbin.Enforcer
}

func NewJwtMiddleware(config *viper.Viper, e *casbin.Enforcer) *JwtMiddleware {
	return &JwtMiddleware{config: config, e: e}
}

func (l *JwtMiddleware) AdminJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(parts[1], l.config)
		if err != nil {
			switch err {
			case utils.ErrExpiredToken:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}
		for k, v := range claims {
			c.Set(k, v)
		}

		if userId, ok := (claims)["user_id"]; ok {
			b, err := l.e.Enforce(utils.ToString(userId), c.Request.URL.Path, "admin")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			if !b {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No permission"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
	// l.e.Enforce()
}
