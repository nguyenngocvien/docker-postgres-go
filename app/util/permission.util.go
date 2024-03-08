package util

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	db "github.com/viennn/docker-postgres-go/app/db/sqlc"
// 	"github.com/viennn/docker-postgres-go/app/token"
// )

// func UpdatePermissionCookie(ctx *gin.Context, username string, permission *map[string]bool) error {
// 	var perms *map[string]bool
// 	if permission == nil {
// 		userInfo, err := db.StoreDB.GetUserInfo(ctx, username)
// 		if err != nil {
// 			log.Println("UpdatePermissionCookie Err : ", err)
// 			return err
// 		}
// 		json.Unmarshal([]byte(userInfo.Permissions), perms)
// 	} else {
// 		perms = permission
// 	}
// 	permissionToken, _, err := token.TokenMaker.CreateToken(username, *perms, Configs.RefreshTokenDuration)
// 	if err != nil {
// 		log.Println("UpdatePermissionCookie Err : ", err)
// 		return err
// 	}
// 	host, secure := SetCookieSameSite(ctx)
// 	ctx.SetCookie("permission", permissionToken, int(Configs.RefreshTokenDuration), "/", host, secure, true)
// 	return nil
// }

// func SetCookieSameSite(ctx *gin.Context) (host string, secure bool) {
// 	mode := Configs.CookieSameSite
// 	var sameSite http.SameSite
// 	switch mode {
// 	case "LAX":
// 		sameSite = http.SameSiteLaxMode
// 	case "STRICT":
// 		sameSite = http.SameSiteStrictMode
// 	case "NONE":
// 		sameSite = http.SameSiteNoneMode
// 	default:
// 		sameSite = http.SameSiteDefaultMode
// 	}
// 	ctx.SetSameSite(sameSite)
// 	host = ctx.Request.Host
// 	secure = Configs.CookieSecure
// 	if !Configs.CookieUseHost {
// 		host = ""
// 	}
// 	return
// }

// func CheckPermission(method []string, typeApi *enums.TypeApi, permissions *map[string]bool) bool {
// 	if method == nil || len(method) == 0 {
// 		return true
// 	}
// 	if len(method) > 1 {
// 		for _, m := range method {
// 			if CheckPermission([]string{m}, typeApi, permissions) {
// 				return true
// 			}
// 		}
// 		return false
// 	}
// 	permissionAPI := fmt.Sprintf("%s-%s", method[0], typeApi)
// 	allPermission := fmt.Sprintf("ALL-%s", typeApi)
// 	if typeApi == nil || (*typeApi) == "" {
// 		return true
// 	}
// 	if permissions == nil {
// 		return false
// 	}
// 	if (*permissions)[permissionAPI] || (*permissions)[allPermission] {
// 		return true
// 	}
// 	fullPerms := "ALL-ALL"
// 	halfAllPerms := fmt.Sprintf("%s-ALL", method[0])
// 	if (*permissions)[fullPerms] || (*permissions)[halfAllPerms] {
// 		return true
// 	}
// 	if method[0] == "GET" {
// 		postPerms := fmt.Sprintf("POST-%s", typeApi)
// 		putPerms := fmt.Sprintf("PUT-%s", typeApi)
// 		deletePerms := fmt.Sprintf("DELETE-%s", typeApi)
// 		return (*permissions)[permissionAPI] || (*permissions)[allPermission] || (*permissions)[postPerms] || (*permissions)[putPerms] || (*permissions)[deletePerms]
// 	}
// 	return false
// }

// func GetPermissionByUser(ctx *gin.Context, username string) (*map[string]bool, error) {
// 	userInfo, err := db.StoreDB.GetUserInfo(ctx, username)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var perms map[string]bool
// 	json.Unmarshal([]byte(userInfo.Permissions), &perms)
// 	return &perms, nil
// }