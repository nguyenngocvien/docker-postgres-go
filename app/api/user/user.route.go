package user

import (
	"github.com/gin-gonic/gin"
	"github.com/viennn/docker-postgres-go/app/middleware"
	"github.com/viennn/docker-postgres-go/app/token"
)

func Routes(routerGroup *gin.RouterGroup) {
	userGroup := routerGroup.Group("/user")
	authUserGroup := userGroup.Use(middleware.AuthMiddleware(token.TokenMaker))
	authUserGroup.POST("/create", createUser)
	authUserGroup.POST("/login", loginUser)
	authUserGroup.POST("/renew_access", renewAccessToken)
}