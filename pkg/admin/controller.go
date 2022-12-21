package admin

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

	routes := r.Group("/api/admin")
	routes.Use(middleware.AdministratorAuthorizationMiddleware)
	routes.GET("/companies/unregistered", h.getUnregisteredCompanies)
	routes.PUT("/companies/:id/verify", h.VerifyCompany)
	routes.DELETE("/companies/:id", h.DeleteCompany)
}
