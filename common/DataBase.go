package common

import (
	"GinProgram/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	//链接信息
	dsn := "root:123456@tcp(127.0.0.1:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"
	//config可自己定义
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	//自动创建数据表
	db.AutoMigrate(&model.User{})

	//将db赋给项目实例
	DB = db

}

func GetDB() *gorm.DB {
	return DB
}
