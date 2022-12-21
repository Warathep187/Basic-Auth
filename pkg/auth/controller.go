package auth

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type Handler struct {
	DB  *mongo.Database
	ctx *context.Context
}

func RegisterRoutes(r *gin.Engine, db *mongo.Database, ctx *context.Context) {
	h := &Handler{
		DB:  db,
		ctx: ctx,
	}

	routes := r.Group("/api/auth")
	routes.POST("/signin", SignInValidator(), h.SignIn)
	routes.POST("/signup", SignUpValidator(), h.SignUp)
	routes.PUT("/verify/:token", h.VerifyAccount)
}
