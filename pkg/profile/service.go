package profile

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/warathep187/BasicAuth/pkg/common/middleware"
	"github.com/warathep187/BasicAuth/pkg/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h Handler) GetProfile(c *gin.Context) {
	user, _ := c.Get("user")

	loggedInUser := user.(middleware.AuthorizedUserPayload)
	userId, _ := primitive.ObjectIDFromHex(loggedInUser.ID)
	if loggedInUser.Role == models.EMPLOYER {
		var company CompanyProfile
		filter := bson.D{{Key: "_id", Value: userId}}
		err := h.DB.Collection("users").FindOne(*h.ctx, filter).Decode(&user)
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200, company)
	} else {
		var user JobSeekerProfile
		filter := bson.D{{Key: "_id", Value: userId}}
		err := h.DB.Collection("users").FindOne(*h.ctx, filter).Decode(&user)
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200, user)
	}
}
