package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	JOB_SEEKER    = "JOB_SEEKER"
	COMPANY       = "COMPANY"
	EMPLOYER      = "EMPLOYER"
	ADMINISTRATOR = "ADMINISTRATOR"
)

type Image struct {
	Key string `json:"key" bson:"key"`
	Url string `json:"url" bson:"url"`
}

type CompanyInformation struct {
	Name        string `json:"name" bson:"name"`
	Information string `json:"information" bson:"information"`
}

type PersonalInformation struct {
	FullName string `json:"fullName" bson:"fullName"`
}

type Contact struct {
	Tel     string `json:"tel" bson:"tel"`
	Email   string `json:"email" bson:"email"`
	Website string `json:"website" bson:"website"`
}

type User struct {
	_id                 primitive.ObjectID
	Email               string              `json:"email" bson:"email"`
	Password            string              `json:"password" bson:"password"`
	Role                string              `json:"role" bson:"role"`
	IsVerified          bool                `json:"isVerified" bson:"isVerified"`
	ProfileImage        Image               `json:"profileImage" bson:"profileImage"`
	CompanyInformation  CompanyInformation  `json:"companyInformation" bson:"companyInformation"`
	PersonalInformation PersonalInformation `json:"personalInformation" bson:"personalInformation"`
	Contact             Contact             `json:"contact" bson:"contact"`
	RegisteredAt        primitive.DateTime  `json:"registeredAt" bson:"registeredAt"`
	VerifiedAt          primitive.DateTime  `json:"verifiedAt" bson:"verifiedAt"`
}
