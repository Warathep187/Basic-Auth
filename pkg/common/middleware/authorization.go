package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/warathep187/BasicAuth/pkg/common/models"
)

type AuthorizedUserPayload struct {
	ID   string
	Role string
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected token")
		}
		return []byte(viper.Get("JWT_AUTHENTICATION_KEY").(string)), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		curTimestamp, _ := strconv.ParseFloat(time.Now().Format("20060102150405"), 64)
		if claims["exp"].(float64) > curTimestamp {
			return nil, fmt.Errorf("Expired token")
		}

		return claims, nil
	} else {
		return nil, err
	}
}

func AllRoleAuthorizationMiddleware(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	token := strings.TrimPrefix(tokenString, "Bearer ")

	decoded, err := ValidateToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	c.Set("user", AuthorizedUserPayload{
		ID:   decoded["sub"].(string),
		Role: decoded["role"].(string),
	})
	c.Next()
}

func AdministratorAuthorizationMiddleware(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	token := strings.TrimPrefix(tokenString, "Bearer ")

	decoded, err := ValidateToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}

	if decoded["role"] != models.ADMINISTRATOR {
		c.JSON(403, gin.H{"message": "Access denied"})
		c.Abort()
		return
	}

	c.Set("user", AuthorizedUserPayload{
		ID:   decoded["sub"].(string),
		Role: decoded["role"].(string),
	})
	c.Next()
}
