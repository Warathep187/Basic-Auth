package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/warathep187/BasicAuth/pkg/admin"
	"github.com/warathep187/BasicAuth/pkg/auth"
	"github.com/warathep187/BasicAuth/pkg/common/db"
	"github.com/warathep187/BasicAuth/pkg/profile"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()
	port := viper.Get("PORT").(string)
	mongodbUri := viper.Get("MONGODB_URI").(string)

	server := gin.Default()
	mongoDB := db.Init(mongodbUri)
	ctx := context.TODO()
	auth.RegisterRoutes(server, mongoDB, &ctx)
	profile.RegisterRoutes(server, mongoDB, &ctx)
	admin.RegisterRoutes(server, mongoDB, &ctx)

	server.Run(port)
}
