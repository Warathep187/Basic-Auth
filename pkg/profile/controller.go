package profile

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/warathep187/BasicAuth/pkg/common/middleware"
	"go.mongodb.org/mongo-driver/mongo"
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

	routes := r.Group("/api/profile")
	routes.GET("/", middleware.AllRoleAuthorizationMiddleware, h.GetProfile)
}
