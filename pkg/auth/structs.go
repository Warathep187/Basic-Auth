package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/warathep187/BasicAuth/pkg/common/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NestedName struct {
	Name string `json:"name" bson:"name"`
}

type NestedFullName struct {
	FullName string `json:"fullName" bson:"fullName"`
}

type UserSignInFindOne struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	Email         string             `json:"email"`
	Password      string             `json:"password"`
	Role          string             `json:"role"`
	ProfileImage  models.Image       `json:"profileImage"`
	JobSeekerName NestedFullName     `json:"jobSeekerName" bson:"personalInformation"`
	CompanyName   NestedName         `json:"companyName" bson:"companyInformation"`
}

type JwtAuthenticationPayload struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type JwtVerificationPayload struct {
	Email string `json:"email"`
}

type JwtAuthenticationCustomClaims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	jwt.StandardClaims
}

type JwtVerificationCustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type SignUpInput struct {
	Type     string `json:"type"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserSignUpFindOne struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Email      string             `json:"email"`
	Role       string             `json:"role"`
	IsVerified bool               `json:"isVerified"`
}

type DecodedVerificationJWT struct {
	Email string `json:"email"`
	Exp   int32  `json:"exp"`
	Iat   int32  `json:"iat"`
	Iss   string `json:"iss"`
}
