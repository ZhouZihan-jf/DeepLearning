package middleware

import (
	"GinProgram/common"
	"GinProgram/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//中间件功能在于存储登录的user信息，从而进一步完成多个验证功能
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		//validate token formate
		//如果token为空或为错误格式，返回权限不足
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		//Bearer部分占七位，故截取7位以后
		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid { //若返回token错误或无效则报错
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort() //Abort是终止运行的意思，直接从调用位置跳出
			return
		}

		//验证通过后获取claim中的userid
		useId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, useId)

		//用户不存在权限不足
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort() //Abort是终止运行的意思，直接从调用位置跳出
			return
		}

		//用户存在，将user信息写入上下文
		ctx.Set("user", user)

		ctx.Next()
	}
}
