package profile

import (
	"github.com/warathep187/BasicAuth/pkg/common/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobSeekerProfile struct {
	ID                  primitive.ObjectID         `json:"id" bson:"_id"`
	Role                string                     `json:"role" bson:"role"`
	ProfileImage        models.Image               `json:"profileImage" bson:"profileImage"`
	Contact             models.Contact             `json:"contact" bson:"contact"`
	PersonalInformation models.PersonalInformation `json:"personalInformation" bson:"personalInformation"`
}

type CompanyProfile struct {
	ID                 primitive.ObjectID        `json:"id" bson:"_id"`
	Role               string                    `json:"role" bson:"role"`
	ProfileImage       models.Image              `json:"profileImage" bson:"profileImage"`
	Contact            models.Contact            `json:"contact" bson:"contact"`
	CompanyInformation models.CompanyInformation `json:"companyInformation" bson:"companyInformation"`
}
