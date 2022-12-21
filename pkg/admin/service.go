package admin

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/warathep187/BasicAuth/pkg/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h Handler) getUnregisteredCompanies(c *gin.Context) {
	filter := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "role", Value: models.EMPLOYER}}, bson.D{{Key: "isVerified", Value: false}}}}}
	sort := options.Find().SetSort(bson.D{{Key: "registeredAt", Value: -1}})

	var companies []*UnregisteredCompany
	cursor, err := h.DB.Collection("users").Find(*h.ctx, filter, sort)

	if err != nil {
		c.JSON(500, gin.H{"message": "Something went wrong"})
		return
	}

	if err = cursor.All(*h.ctx, &companies); err != nil {
		panic(err)
	}
	for _, result := range companies {
		json.Marshal(result)
	}

	c.JSON(200, gin.H{"companies": companies})
}

func (h Handler) VerifyCompany(c *gin.Context) {
	companyId := c.Param("id")
	if companyId == "" {
		c.JSON(400, gin.H{"message": "Company ID must be provided"})
		return
	}
	companyObjectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(400, gin.H{"message": "ID is invalid"})
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isVerified", Value: true}}}}

	result, err := h.DB.Collection("users").UpdateByID(*h.ctx, companyObjectId, update)
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not verify this account"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(409, gin.H{"message": "Company not found"})
		return
	}
	c.JSON(204, gin.H{})
}

func (h Handler) DeleteCompany(c *gin.Context) {
	companyId := c.Param("id")
	if companyId == "" {
		c.JSON(400, gin.H{"message": "Company ID must be provided"})
		return
	}
	companyObjectId, err := primitive.ObjectIDFromHex(companyId)
	if err != nil {
		c.JSON(400, gin.H{"message": "ID is invalid"})
	}
	filter := bson.D{{Key: "_id", Value: companyObjectId}}
	_, err = h.DB.Collection("users").DeleteOne(*h.ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not remove company, Please try again"})
		return
	}
	c.JSON(204, gin.H{})
}
