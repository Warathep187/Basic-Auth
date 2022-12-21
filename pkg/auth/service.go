package auth

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/warathep187/BasicAuth/pkg/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var sess *session.Session
var sesService *ses.SES

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})
	if err != nil {
		panic(err)
	}

	sesService = ses.New(sess)

	fmt.Println("Hello, AWS")
}

// Helpers
func GenerateAuthenticationJWT(payload *JwtAuthenticationPayload) (string, error) {
	key := []byte(viper.Get("JWT_AUTHENTICATION_KEY").(string))
	claims := &JwtAuthenticationCustomClaims{
		payload.ID,
		payload.Role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3 * 24).Unix(),
			Issuer:    "authentication",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}

func GenerateVerificationJWT(payload *JwtVerificationPayload) (string, error) {
	key := []byte(viper.Get("JWT_VERIFICATION_KEY").(string))
	claims := &JwtVerificationCustomClaims{
		payload.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
			Issuer:    "verification",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func SendVerificationEmail(sendTo string) {
	sender := viper.Get("SENDER_EMAIL").(string)

	token, err := GenerateVerificationJWT(&JwtVerificationPayload{Email: sendTo})
	if err != nil {
		panic(err)
	}

	htmlBody := fmt.Sprintf(`
		<h1>Verification Email</h1>
		<a href='%s'>Click Here To Verify your Account</a>
 	`, viper.Get("CLIENT_URL").(string)+"/verify/"+token)

	inp := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(sendTo),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String("Hello, world"),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Verification Procedure"),
			},
		},
		Source: aws.String(sender),
	}

	result, err := sesService.SendEmail(inp)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent!!")
	fmt.Println(result)
}

func ValidateVerificationToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected headers")
		}
		return []byte(viper.Get("JWT_VERIFICATION_KEY").(string)), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["exp"].(float64) > claims["iat"].(float64) {
			return nil, fmt.Errorf("Expired token")
		}

		return claims, nil
	} else {
		return nil, err
	}
}

// Main
func (h Handler) SignIn(c *gin.Context) {
	body, _ := c.Get("body")
	signinInput := body.(SignInInput)

	var user *UserSignInFindOne
	filter := bson.D{{Key: "email", Value: signinInput.Email}, {Key: "isVerified", Value: true}}
	h.DB.Collection("users").FindOne(*h.ctx, filter).Decode(&user)
	if user == nil {
		c.JSON(400, gin.H{"message": "Account not found"})
		return
	}
	err := ComparePassword(user.Password, signinInput.Password)
	if err != nil {
		c.JSON(400, gin.H{"message": "Password is incorrect"})
		return
	}

	token, err := GenerateAuthenticationJWT(&JwtAuthenticationPayload{ID: user.ID.Hex(), Role: user.Role})
	if err != nil {
		c.JSON(500, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(200, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h Handler) SignUp(c *gin.Context) {
	body, _ := c.Get("body")

	signUpInput := body.(SignUpInput)
	c.BindJSON(signUpInput)

	filter := bson.D{bson.E{Key: "email", Value: signUpInput.Email}}
	var user *UserSignUpFindOne
	h.DB.Collection("users").FindOne(*h.ctx, filter).Decode(&user)

	var role string
	if signUpInput.Type == models.COMPANY {
		role = models.EMPLOYER
	} else {
		role = models.JOB_SEEKER
	}
	if user == nil {
		hashedPassword, err := HashPassword(signUpInput.Password)
		if err != nil {
			c.JSON(500, gin.H{"message": "Something went wrong"})
			return
		}
		var document bson.D
		if role == models.EMPLOYER {
			document = bson.D{
				{Key: "email", Value: signUpInput.Email},
				{Key: "password", Value: hashedPassword},
				{Key: "role", Value: role},
				{Key: "isVerified", Value: false},
				{Key: "profileImage", Value: bson.D{{Key: "key", Value: ""}, {Key: "url", Value: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"}}},
				{Key: "companyInformation", Value: bson.D{{Key: "name", Value: signUpInput.Name}, {Key: "information", Value: ""}, {Key: "images", Value: []interface{}{}}}},
				{Key: "contact", Value: bson.D{{Key: "tel", Value: ""}, {Key: "email", Value: ""}, {Key: "website", Value: ""}}},
				{Key: "registeredAt", Value: time.Now()},
				{Key: "verifiedAt", Value: nil},
			}
			_, err := h.DB.Collection("users").InsertOne(*h.ctx, document)
			if err != nil {
				c.JSON(500, gin.H{"message": "Something went wrong"})
				return
			}
			c.JSON(201, gin.H{"message": "Please wait for administrator confirm your company account"})
		} else {
			document = bson.D{
				{Key: "email", Value: signUpInput.Email},
				{Key: "password", Value: hashedPassword},
				{Key: "role", Value: role},
				{Key: "isVerified", Value: false},
				{Key: "profileImage", Value: bson.D{{Key: "key", Value: ""}, {Key: "url", Value: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"}}},
				{Key: "contact", Value: bson.D{{Key: "tel", Value: ""}, {Key: "email", Value: ""}, {Key: "website", Value: ""}}},
				{Key: "registeredAt", Value: time.Now()},
				{Key: "personalInformation", Value: bson.D{{Key: "fullName", Value: signUpInput.Name}, {Key: "tags", Value: []interface{}{}}}},
				{Key: "favoriteJobs", Value: []interface{}{}},
				{Key: "verifiedAt", Value: nil},
			}
			_, err := h.DB.Collection("users").InsertOne(*h.ctx, document)
			if err != nil {
				c.JSON(500, gin.H{"message": "Something went wrong"})
				return
			}
			go SendVerificationEmail(signUpInput.Email)
			c.JSON(201, gin.H{"message": "Please verify your email"})
		}
	} else {
		if user.IsVerified {
			c.JSON(409, gin.H{"message": "Email has already used"})
			return
		}
		if user.Role == models.EMPLOYER && signUpInput.Type == models.JOB_SEEKER {
			c.JSON(409, gin.H{"message": "Email has already used"})
			return
		}
		if user.Role == models.JOB_SEEKER && signUpInput.Type == models.COMPANY {
			c.JSON(409, gin.H{"message": "Email has already used"})
			return
		}
		var update bson.D

		hashedPassword, err := HashPassword(signUpInput.Password)
		if err != nil {
			c.JSON(500, gin.H{"message": "Something went wrong"})
		}
		if role == models.EMPLOYER {
			update = bson.D{
				{
					Key: "$set",
					Value: bson.D{
						{Key: "password", Value: hashedPassword},
						{Key: "companyInformation", Value: bson.D{{Key: "name", Value: signUpInput.Name}, {Key: "information", Value: ""}, {Key: "images", Value: []interface{}{}}}},
					},
				},
			}
			_, err := h.DB.Collection("users").UpdateByID(*h.ctx, user.ID, update)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{"message": "Something went wrong"})
				return
			}
			c.JSON(201, gin.H{"message": "Please wait for administrator confirm your company account"})
		} else {
			update = bson.D{
				{
					Key: "$set",
					Value: bson.D{
						{Key: "password", Value: hashedPassword},
						{Key: "personalInformation", Value: bson.D{{Key: "fullName", Value: signUpInput.Name}, {Key: "tags", Value: []interface{}{}}}},
					},
				},
			}
			_, err := h.DB.Collection("users").UpdateByID(*h.ctx, user.ID, update)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{"message": "Something went wrong"})
				return
			}
			go SendVerificationEmail(signUpInput.Email)
			c.JSON(201, gin.H{"message": "Please verify your email"})
		}
	}
}

func (h Handler) VerifyAccount(c *gin.Context) {
	token := c.Param("token")

	result, err := ValidateVerificationToken(token)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	email := result["email"].(string)

	filter := bson.D{{Key: "$and", Value: bson.A{bson.D{{Key: "email", Value: email}}, bson.D{{Key: "isVerified", Value: false}}}}}
	count, err := h.DB.Collection("users").CountDocuments(*h.ctx, filter)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"message": "Something went wrong"})
		return
	}
	if count == 0 {
		c.JSON(409, gin.H{"message": "Account is verified"})
		return
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isVerified", Value: true}, {Key: "verifiedAt", Value: time.Now()}}}}
	_, err = h.DB.Collection("users").UpdateOne(*h.ctx, filter, update)
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not verify your account"})
		return
	}

	c.JSON(200, gin.H{"message": "Verified"})
}
