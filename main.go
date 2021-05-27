package main

import (
	"GinProgram/common"

	"github.com/gin-gonic/gin"
)

func main() {
	//链接数据库
	common.InitDb()
	/*
		db := common.GetDB()
		defer db.Close()//现版本无close方法
	*/
	r := gin.Default()
	//r调用方法
	r = CollectRoute(r)
	//main内运行
	r.Run() // listen and serve on 0.0.0.0:8080
	panic(r.Run())
}
