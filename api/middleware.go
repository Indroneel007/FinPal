package api

import (
	"examples/SimpleBankProject/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	authorizationHeaderKey        = "Authorization"
	authorizationHeaderBearerType = "bearer"
	authorizationPayloadKey       = "authorization_payload"
)

func AuthMiddleware(tokenMaker util.Maker) gin.HandlerFunc {
	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}

	viper.AutomaticEnv()

	return func(c *gin.Context) {
		duration := viper.GetDuration("TOKEN_DURATION")
		if duration == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "TOKEN_DURATION is not set in the environment variables"})
			return
		}

		authorizationHeader := c.GetHeader(authorizationHeaderKey)

		if authorizationHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != authorizationHeaderBearerType {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token := fields[1]
		payload, err := tokenMaker.VerifyToken(token, duration)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}
