package admin

import (
	"github.com/warathep187/BasicAuth/pkg/common/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnregisteredCompany struct {
	ID                 primitive.ObjectID        `json:"id" bson:"_id"`
	Email              string                    `json:"email" bson:"email"`
	Role               string                    `json:"role" bson:"role"`
	ProfileImage       models.Image              `json:"profileImage" bson:"profileImage"`
	Contact            models.Contact            `json:"contact" bson:"contact"`
	CompanyInformation models.CompanyInformation `json:"companyInformation" bson:"companyInformation"`
	RegisteredAt       primitive.DateTime        `json:"registeredAt" bson:"registeredAt"`
}
