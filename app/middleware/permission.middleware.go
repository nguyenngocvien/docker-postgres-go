package middleware

// import (
// 	"errors"

// 	"github.com/gin-gonic/gin"
// 	"github.com/viennn/docker-postgres-go/app/util"
// 	"github.com/viennn/docker-postgres-go/app/util/enums"
// )

// func PermissionMiddleware(typeApi enums.TypeApi) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		util.SetCookieSameSite(ctx)
// 	cookie, err :=ctx.Cookie("permission")
// 	if util.Configs.DefaultAuthenticationUsername!="" && err !=nil{
// 		err=nil
// 		perms,err:=util.GetPermissionByUser(ctx,util.Configs.DefaultAuthenticationUsername)
// 		if err != nil {
// 			ctx.AbortWithStatusJSON(403, util.ErrorResponse(err))
// 			return
// 		}
// 		cookie,_,err = token.TokenMaker.CreateToken(util.Configs.DefaultAuthenticationUsername,*perms,util.Configs.RefreshTokenDuration)
// 	}
// 	if err != nil {
// 		ctx.AbortWithStatusJSON(403, util.ErrorResponse(err))
// 		return
// 	}
// 	if cookie == ""{
// 		ctx.AbortWithStatusJSON(403, util.ErrorResponse(errors.New("user have not permission or permission expired, please refresh (F5) page or contact admin")))
// 		return
// 	}
// 	payload,err:=token.TokenMaker.VerifyToken(cookie)
// 	if err != nil {
// 		ctx.AbortWithStatusJSON(403, util.ErrorResponse(err))
// 		return
// 	}
// 	isValid := util.CheckPermission([]string{ctx.Request.Method},&typeApi,&payload.Permission)
// 	if !isValid {
// 		ctx.AbortWithStatusJSON(403, util.ErrorResponse(errors.New("Tài khoản của bạn không có quyền truy cập vào chức năng này")))
// 		return
// 	}
// 		ctx.Next()
// 	}
// }