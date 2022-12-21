package auth

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/warathep187/BasicAuth/pkg/common/models"
)

func SignInValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var signInInput SignInInput
		if err := c.ShouldBindJSON(&signInInput); err != nil {
			c.JSON(400, gin.H{"message": "Invalid input format"})
			c.Abort()
			return
		}
		if strings.Trim(signInInput.Email, " ") == "" {
			c.JSON(400, gin.H{"message": "Email must be provided"})
			c.Abort()
		} else if matched, _ := regexp.Match(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, []byte(signInInput.Email)); !matched {
			c.JSON(400, gin.H{"message": "Email is invalid"})
			c.Abort()
		} else if len(strings.Trim(signInInput.Password, " ")) < 6 {
			c.JSON(400, gin.H{"message": "Password must be at least 6 characters"})
			c.Abort()
		} else {
			c.Set("body", signInInput)
			c.Next()
		}
	}
}

func SignUpValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var signUpInput SignUpInput

		if err := c.ShouldBindJSON(&signUpInput); err != nil {
			c.JSON(400, gin.H{"message": "Invalid input format"})
			c.Abort()
			return
		}
		if signUpInput.Type != models.JOB_SEEKER && signUpInput.Type != models.COMPANY {
			c.JSON(400, gin.H{"message": "Invalid sign-up type"})
			c.Abort()
		} else if matched, _ := regexp.Match(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, []byte(signUpInput.Email)); !matched {
			c.JSON(400, gin.H{"message": "Email is invalid"})
			c.Abort()
		} else if len(strings.Trim(signUpInput.Password, " ")) < 6 {
			c.JSON(400, gin.H{"message": "Password must be at least 6 characters"})
			c.Abort()
		} else if len(strings.Trim(signUpInput.Name, " ")) == 0 {
			c.JSON(400, gin.H{"message": "Name must be provided"})
			c.Abort()
		} else {
			c.Set("body", signUpInput)
			c.Next()
		}
	}
}
