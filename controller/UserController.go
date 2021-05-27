package controller

import (
	"GinProgram/common"
	"GinProgram/model"
	"log"
	"net/http"

	"GinProgram/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context) {
	//获取数据库
	DB := common.GetDB()

	//获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}
	//如果name为空，随机给10位的字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Println(name, telephone, password)
	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 { //若查询结果为空则ID为0
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "该号码已经存在"})
		return
	}
	//创建用户
	//给密码加密
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "加密失败"})
		return
	}
	nweUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&nweUser)
	//返回结果
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	// 获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}
	//判断存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 { //若查询结果为空则ID为0
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//判断密码正误
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}
	//发放token
	//token分三部分：加密协议头.荷载内容.前两部分加key再hash(作为签名)
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error : %v", err)
		return
	}

	//返回结果
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": user}})
}
