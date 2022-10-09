package controller

import (
	"github.com/gin-gonic/gin"
	"tc-back/utils"
)

func InitRouter(r *gin.Engine) {
	r.NoRoute(utils.Default404Router)
	r.Static("/", utils.FrontCfg.Root)

	v1 := r.Group("api")
	{
		//注册
		v1.POST("/reg", Reg)
		//登陆，返回token
		v1.POST("/login", Login)
	}
	//前端没有把token放head，这里手动校验，不用middleware了
	//v1.Use(middleware.JWTAuth)
	{
		//md5秒传
		v1.POST("/md5", MD5)
		//真实上传
		v1.POST("/upload", UploadFile)
		//获取用户文件信息
		v1.POST("/myfiles", MyFiles)
		//获取共享列表 获取共享文件信息
		v1.POST("/sharefiles", ShareFiles)
		//文件分享 文件删除 文件下载pv++
		v1.POST("/dealfile", DealFile)
		//取消共享文件 转存共享文件 共享文件pv下载+1
		v1.POST("/dealsharefile", DealShareFile)
		//图片分享 请求浏览  图片信息 取消分享
		v1.POST("/sharepic", SharePic)
	}
}
