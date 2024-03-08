package account

import (
	"github.com/gin-gonic/gin"
	"github.com/viennn/docker-postgres-go/app/middleware"
	"github.com/viennn/docker-postgres-go/app/token"
)

func Routes(routerGroup *gin.RouterGroup) {
	userGroup := routerGroup.Group("/accounts")
	authUserGroup := userGroup.Use(middleware.AuthMiddleware(token.TokenMaker))
	authUserGroup.POST("/create", createAccount)
	authUserGroup.GET("/:id", getAccount)
	authUserGroup.GET("/accounts", listAccounts)
}