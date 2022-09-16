package controller

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"tc-back/dfst"
	"tc-back/middleware"
	"tc-back/opdb"
	"tc-back/oprds"
	"tc-back/utils"
	"time"
)

//注册接口的回调函数
func Reg(c *gin.Context) {
	var userInfo dfst.RegisterInfoReq
	//将body中的json内容反射到结构体中
	bindErr := c.BindJSON(&userInfo)
	//如果错误说明client那边发来的数据不对，直接报错即可
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	//将用户信息插入数据库中
	err := opdb.InsertRecordToUserInfo(userInfo)
	//如果出错说明该用户已经注册过了
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

//登陆接口的回调函数
func Login(c *gin.Context) {
	var loginReq dfst.LoginInfoReq
	//将body中的json内容反射到结构体中
	bindErr := c.BindJSON(&loginReq)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  1,
			"token": "faild",
		})
		return
	}
	//验证账号密码是否正确，即是否在数据库中存在，注册过
	pass, _ := opdb.HasFileInUserInfo(loginReq)
	if !pass {
		c.JSON(http.StatusOK, gin.H{
			"code":  1,
			"token": "faild",
		})
		return
	}

	//创建一个token GenerateToken
	token, err := middleware.GenerateToken(c, loginReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  1,
			"token": "faild",
		})
		return
	}
	//token 存数据库
	oprds.SetLoginToken(loginReq.User, token)
	//将对象名和token返回
	c.JSON(http.StatusOK, gin.H{
		"code":  "0",
		"token": token,
	})
}

// 1. 查询数据库是否存在该文件，如果不存在则返回
// 2. 如果存在，查询是否是自己上传的，如果是则返回
// 3. 如果不是自己上传的
//       文件的引用+1
//  	 更新用户的文件列表和用户文件数量集

func MD5(c *gin.Context) {
	var md5req dfst.MD5Req
	//将body中的json内容反射到结构体中
	bindErr := c.BindJSON(&md5req)
	if bindErr != nil {
		log.Println("bind")

		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	//token 校验
	if err := oprds.CheckToekn(md5req.User, md5req.Token); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 3,
		})
		return
	}
	//1. 数据库查询是否有对应的md5文件
	has := opdb.HasFileInFileInFoByMd5(md5req.Md5)
	//如果有
	if has {
		//看看是不是用户自己上传的，如果是则无需重复上传
		hashas := opdb.HasFileInUserFileList(md5req.User, md5req.Md5, md5req.Filename)

		//说明此用户已经保存此文件
		if hashas {
			log.Println("have")

			c.JSON(http.StatusOK, gin.H{
				"code": 2,
			})
			return
		} else {
			//说明不是自己上传的,将文件的引用计数+1
			err := opdb.UpdateFileInfoCountInc(md5req.Md5)
			if err != nil {
				log.Println("UpdateFileInfoCountInc")

				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
			//用户文件列表插入一条数据
			err = opdb.InsertRecordToUserFileList(md5req.User, md5req.Filename, md5req.Md5)
			if err != nil {
				log.Println("InsertRecordToUserFileList")

				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
			//更新用户文件数量集，如果没有则插入一条数据
			err = opdb.UpdateUserFileListCountInc(md5req.User)
			if err != nil {
				log.Println("UpdateUserFileListCountInc")

				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
			// return
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
			})
			return
		}
	} else {
		//没有文件，妙传失败
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
	}
}

func UploadFile(c *gin.Context) {
	//解析multipart/form-data
	user := c.PostForm("user")
	md5_ := c.PostForm("md5")
	size, _ := strconv.Atoi(c.PostForm("size"))
	file, _ := c.FormFile("file")

	ext := strings.Split(file.Filename, ".")
	fileExt := ext[len(ext)-1]

	f, _ := file.Open()
	defer f.Close()

	buf := &bytes.Buffer{}
	_, _ = buf.ReadFrom(f)
	//存储到fastdfs中
	fileID, err := utils.UploadFileByBuffer(file.Filename, fileExt, buf.Bytes())
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	//更新数据库相关信息
	url := utils.MakeFileURL(fileID)

	//插入一条信息 到 文件信息表
	err = opdb.InsertRecordToFileInfo(md5_, fileID, url, fileExt, size, 1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	//用户文件列表插入一条数据
	err = opdb.InsertRecordToUserFileList(user, file.Filename, md5_)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	//更新用户文件数量集，如果没有则插入一条数据
	err = opdb.UpdateUserFileListCountInc(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}
func MyFiles(c *gin.Context) {
	var myFileReq dfst.MyFilesReq
	//将body中的json内容反射到结构体中
	bindErr := c.ShouldBind(&myFileReq)
	if bindErr != nil {
		log.Println("bindErr err", bindErr)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	//token 校验
	if err := oprds.CheckToekn(myFileReq.User, myFileReq.Token); err != nil {
		log.Println("token err:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	cmd := c.Query("cmd")
	if cmd == "count" { //获取用户文件数量
		tot, err := opdb.GetUserFileCount(myFileReq.User)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":  0,
				"total": tot,
			})
			return
		}
		//获取用户文件信息 normal
		//按下载量升序 pvasc
		//按下载量降序 pvdesc
	} else if cmd == "normal" || cmd == "pvasc" || cmd == "pvdesc" {
		if myFileReq.Count == 0 {
			myFileReq.Count = 10
		}
		fileInfos, err := opdb.GetUserFileListJoinFileInfo(myFileReq.User, cmd, myFileReq.Start, myFileReq.Count)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		tot, err := opdb.GetUserFileCount(myFileReq.User)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": len(fileInfos),
			"total": tot,
			"files": fileInfos,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}

}

//获取共享文件个数 /api/sharefiles?cmd=count
//获取共享文件列表 /api/sharefiles?cmd=normal
//获取共享文件下载排行榜 /api/sharefiles?cmd=pvdesc/pvasc

func ShareFiles(c *gin.Context) {
	//获取共享文件个数
	cmd := c.Query("cmd")
	if cmd == "count" {
		//获取共享文件个数
		tot, err := opdb.GetUserFileCount("xxx_share_xxx_file_xxx_list_xxx_count_xxx")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"total": tot,
		})
		return
	}

	var shareFileReq dfst.ShareFilesReq
	//将body中的json内容反射到结构体中
	bindErr := c.BindJSON(&shareFileReq)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}

	if cmd == "normal" {
		//获取共享文件列表  normal
		fileInfos, err := opdb.GetShareFileListJoinFileInfo(shareFileReq.Start, shareFileReq.Count)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		tot, err := opdb.GetUserFileCount("xxx_share_xxx_file_xxx_list_xxx_count_xxx")
		if err != nil {
			fmt.Println("GetUserFileCount err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		for i := range fileInfos {
			fileInfos[i].SharedStatus = 1
		}

		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": len(fileInfos),
			"total": tot,
			"files": fileInfos,
		})
	} else {
		//获取共享文件下载排行榜  pvdesc
		tot, err := opdb.GetUserFileCount("xxx_share_xxx_file_xxx_list_xxx_count_xxx")
		if err != nil {
			fmt.Println("GetUserFileCount err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}

		//TODO 共享文件下载排行榜 redis
		mysqlNum := tot
		redisNum := oprds.GetShareFileNum("FILE_PUBLIC_ZSET")
		//mysql共享文件数量和redis共享文件数量对比，判断是否相等
		if mysqlNum != redisNum {
			//如果不相等，清空redis数据，重新从mysql中导入数据到redis (mysql和redis交互)
			oprds.DelKey("FILE_PUBLIC_ZSET")
			oprds.DelKey("FILE_NAME_HASH")
			//从mysql中导入数据到redis
			fileInfos, err := opdb.GetShareFileList()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
			for _, info := range fileInfos {
				value := fmt.Sprintf("%s%s", info.Md5, info.FileName)
				score := info.Pv
				//增加有序集合成员
				oprds.SetZsetKey("FILE_PUBLIC_ZSET", score, value)
				//增加hash记录
				oprds.SetHashKey("FILE_NAME_HASH", value, info.FileName)
			}
		}
		//现在redis和mysql的数据是同步的
		//降序获取有序集合的元素
		md5filename, err := oprds.GetZsetZrevrange("FILE_PUBLIC_ZSET", 0, tot)
		if err != nil {
			log.Println("redis err zset zrevrange ", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		var resp []dfst.ShareFileList2Resp
		for _, md2f := range md5filename {
			//filename
			filename, err := oprds.GetHashFilename("FILE_NAME_HASH", md2f)
			if err != nil {
				log.Println("redis err zset zrevrange ", err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
			temp := dfst.ShareFileList2Resp{
				FileName: filename,
			}
			//pv
			pv, err := oprds.GetZsetScore("FILE_PUBLIC_ZSET", md2f)
			if err != nil {
				log.Println("redis err zset zrevrange ", err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
			temp.Pv = int(pv)
			resp = append(resp, temp)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": len(resp),
			"total": tot,
			"files": resp,
		})
	}
}
func DealFile(c *gin.Context) {
	var delfileReq dfst.DealFileReq
	//将body中的json内容反射到结构体中
	bindErr := c.BindJSON(&delfileReq)
	if bindErr != nil {
		log.Println("bind json err")
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	//token 校验
	if err := oprds.CheckToekn(delfileReq.User, delfileReq.Token); err != nil {
		log.Println("token err")
		c.JSON(http.StatusOK, gin.H{
			"code": 4,
		})
		return
	}
	//cmd
	cmd := c.Query("cmd")
	//cmd=share del pv
	//文件分享 文件删除 文件下载pv++
	fileid := fmt.Sprintf("%s%s", delfileReq.Md5, delfileReq.Filename)
	if cmd == "share" {
		has, err := oprds.CheckZsetHasFile("FILE_PUBLIC_ZSET", fileid)
		if err != nil {
			log.Println("redis err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		if has { //存在直接返回
			log.Println("has")

			c.JSON(http.StatusOK, gin.H{
				"code": 3,
			})
			return
		}
		//更新共享文件标志
		err = opdb.UpdateUserFileListForShare(delfileReq.User, delfileReq.Md5, delfileReq.Filename)
		if err != nil {
			log.Println("UpdateUserFileListForShare ", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//分享文件的信息，额外保存在share_file_list保存列表
		err = opdb.InsertRecordToShareFileList(delfileReq.User, delfileReq.Md5, delfileReq.Filename)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//查询共享文件的数量 不存在->创建，存在-> +1
		err = opdb.UpdateUserFileListCountInc("xxx_share_xxx_file_xxx_list_xxx_count_xxx")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//存入redis
		//增加有序集合成员
		oprds.SetZsetKey("FILE_PUBLIC_ZSET", 0, fileid)
		//增加hash记录
		oprds.SetHashKey("FILE_NAME_HASH", fileid, delfileReq.Filename)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	} else if cmd == "del" {
		//删除文件
		has, err := oprds.CheckZsetHasFile("FILE_PUBLIC_ZSET", fileid)
		if err != nil {
			log.Println("redis err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		if has { //说明共享了
			//删除分享列表(share_file_list)的数据
			err = opdb.DelRecordForShareFileList(delfileReq.User, delfileReq.Md5, delfileReq.Filename)
			if err != nil {
				log.Println("DelRecordForShareFileList err", err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
			//共享文件的数量-1
			err = opdb.UpdateUserFileCountDecr("xxx_share_xxx_file_xxx_list_xxx_count_xxx")
			if err != nil {
				log.Println("UpdateUserFileCountDecr err", err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
			//删除redis记录
			oprds.DelZsetKey("FILE_PUBLIC_ZSET", fileid)
			oprds.DelHashKey("FILE_NAME_HASH", fileid)
		}
		//用户文件数量-1  user_file_count
		err = opdb.UpdateUserFileCountDecr(delfileReq.User)
		if err != nil {
			log.Println("UpdateUserFileCountDecr err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//删除用户文件列表数据
		err = opdb.DelRecordForUserFileList(delfileReq.User, delfileReq.Md5, delfileReq.Filename)
		if err != nil {
			log.Println("UpdateUserFileCountDecr err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		count, err := opdb.GetFileInfoCount(delfileReq.Md5)
		if err != nil {
			log.Println("GetFileInfoCount err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		if count > 0 {
			//文件信息表(file_info)的文件引用计数count，减1
			err = opdb.UpdateFileInfoCountDecr(delfileReq.Md5)
			if err != nil {
				log.Println("UpdateFileInfoCountDecr err", err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
		}
		count--
		//说明没有用户引用此文件，需要在storage删除此文件
		if count == 0 {
			//查询文件的id
			storagefileld, err := opdb.GetFileInfoFileId(delfileReq.Md5)
			if err != nil {
				log.Println("GetFileInfoFileId err", err)
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
				})
				return
			}
			//删除文件信息表中该文件的信息
			opdb.DelRecordForFileInfo(delfileReq.Md5)

			//从storage服务器删除此文件，参数为为文件id
			_ = utils.DelFileByFileID(storagefileld)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	} else if cmd == "pv" {
		//更新文件下载计数
		opdb.UpdateUserFileListPvInc(delfileReq.User, delfileReq.Md5, delfileReq.Filename)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
}

// 取消分享文件 cancel
// 转存文件 save
// 共享文件下载pv pv

func DealShareFile(c *gin.Context) {
	var dealShareFileReq dfst.DealShareFileReq
	//将body中的json内容反射到结构体中
	bindErr := c.BindJSON(&dealShareFileReq)
	if bindErr != nil {
		log.Println("bind json err")
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	//cmd
	cmd := c.Query("cmd")

	if cmd == "cancel" {
		// 取消分享文件 cancel

		// 1. user_file_list对应shared_status置0
		err := opdb.UpdateUserFileListForNotShare(dealShareFileReq.User, dealShareFileReq.Md5, dealShareFileReq.Filename)
		if err != nil {
			log.Println("UpdateUserFileListForNotShare ", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		// 2. user_file_count的xxx_share_xxx_file_xxx_list_xxx_count_xxx减1
		//共享文件的数量-1
		err = opdb.UpdateUserFileCountDecr("xxx_share_xxx_file_xxx_list_xxx_count_xxx")
		if err != nil {
			log.Println("UpdateUserFileCountDecr err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		// 3. share_file_list对应的信息删除
		//删除分享列表(share_file_list)的数据
		err = opdb.DelRecordForShareFileList(dealShareFileReq.User, dealShareFileReq.Md5, dealShareFileReq.Filename)
		if err != nil {
			log.Println("DelRecordForShareFileList err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		// 4. redis删除
		//删除redis记录
		fileid := dealShareFileReq.Md5 + dealShareFileReq.Filename
		oprds.DelZsetKey("FILE_PUBLIC_ZSET", fileid)
		oprds.DelHashKey("FILE_NAME_HASH", fileid)
		// 5. return
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	} else if cmd == "save" {
		// 转存文件 save
		//1.先查询是个人文件列表是否已经存在该文件。如果存在则返回5
		has := opdb.HasFileInUserFileList(dealShareFileReq.User, dealShareFileReq.Md5, dealShareFileReq.Filename)
		if has {
			c.JSON(http.StatusOK, gin.H{
				"code": 5,
			})
			return
		}
		//2.增加 file_info 表的 count 计数，表示多一个人保存了该文件。
		err := opdb.UpdateFileInfoCountInc(dealShareFileReq.Md5)
		if err != nil {
			log.Println("DelRecordForShareFileList err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//3.个人的 user_file_list 增加一条文件记录
		err = opdb.InsertRecordToUserFileList(dealShareFileReq.User, dealShareFileReq.Filename, dealShareFileReq.Md5)
		if err != nil {
			log.Println("DelRecordForShareFileList err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//4.更新个人的 user_file_count
		err = opdb.UpdateUserFileListCountInc(dealShareFileReq.User)
		if err != nil {
			log.Println("DelRecordForShareFileList err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//5.return
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return

	} else if cmd == "pv" {
		// 共享文件下载pv pv
		//1. 更新 share_file_list 的 pv 值
		opdb.UpdateShareFileListPvInc(dealShareFileReq.Md5, dealShareFileReq.Filename)
		//2. 更新 redis 里的 FILE_PUBLIC_ZSET，用作排行榜
		fileid := dealShareFileReq.Md5 + dealShareFileReq.Filename
		oprds.IncZsetKey("FILE_PUBLIC_ZSET", fileid)
		//3. return
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
}

// 图片分享 share
// 图片浏览 browse
// 图片信息 normal
// 取消分享 cancel

func SharePic(c *gin.Context) {
	var sharePicReq dfst.SharePicReq
	//将body中的json内容反射到结构体中
	bindErr := c.BindJSON(&sharePicReq)
	if bindErr != nil {
		log.Println("ShouldBind json err")
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
	if sharePicReq.Token != "" {
		//token 校验
		if err := oprds.CheckToekn(sharePicReq.User, sharePicReq.Token); err != nil {
			log.Println("token err")

			c.JSON(http.StatusOK, gin.H{
				"code": 4,
			})
			return
		}
	}
	//cmd
	cmd := c.Query("cmd")
	if cmd == "share" {
		// 1. 生成要返回的url md5（时间+随机数）
		rand.Seed(time.Now().UnixNano())
		randomNum := rand.Intn(8999) + 1000 // 生成随机数
		tmp := fmt.Sprintf("%s%d%d", sharePicReq.Filename, time.Now().UnixNano(), randomNum)
		urlmd51 := md5.Sum([]byte(tmp))
		urlmd5 := base64.StdEncoding.EncodeToString(urlmd51[:])
		// 2. share_picture_list添加相应信息
		err := opdb.InsertRecordToSharePicList(sharePicReq.User, sharePicReq.Md5, sharePicReq.Filename, urlmd5)
		if err != nil {
			log.Println("InsertRecordToSharePicList err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		// 3. 共享图片增加 _share_picture_list_count
		err = opdb.UpdateUserFileListCountInc(sharePicReq.User + "_share_picture_list_count")
		if err != nil {
			log.Println("UpdateUserFileListCountInc err", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		// 4. return
		c.JSON(http.StatusOK, gin.H{
			"code":   0,
			"urlmd5": urlmd5,
		})
		return

	} else if cmd == "normal" {
		//获取共享文件列表
		// 1. 获取相关信息
		sharePicInfos, err := opdb.GetSharePicListJoinFileInfo(sharePicReq.User, sharePicReq.Start, sharePicReq.Count)
		if err != nil {
			log.Println(sharePicInfos)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		count := len(sharePicInfos)
		tot, err := opdb.GetUserFileCount(sharePicReq.User)
		if err != nil {
			log.Println(sharePicInfos)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		// 2. return
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": count,
			"total": tot,
			"files": sharePicInfos,
		})
		return

	} else if cmd == "browse" {
		//1. share_picture_list
		picInfo, err := opdb.GetSharePicListInfo(sharePicReq.UrlMd5)
		if err != nil {
			log.Println(picInfo, err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//2. file_info
		url, err := opdb.GetFileInfoUrl(picInfo.FileMd5)
		if err != nil {
			log.Println("GetFileInfoUrl", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//3. update share_picture_list.pv
		opdb.UpdateSharePicListPvInc(sharePicReq.UrlMd5)
		//4. return
		resp := dfst.SharePicListResp{
			Code:       0,
			User:       picInfo.User,
			Pv:         picInfo.Pv + 1,
			CreateTime: picInfo.CreateTime,
			Url:        url,
		}
		c.JSON(http.StatusOK, resp)
		return

	} else if cmd == "cancel" {
		//1. share_picture_list 删除相应信息
		err := opdb.DelRecordForSharePicList(sharePicReq.User, sharePicReq.UrlMd5)
		if err != nil {
			log.Println("GetFileInfoUrl", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//2.  _share_picture_list_count 共享图片减1
		err = opdb.UpdateUserFileCountDecr(sharePicReq.User + "_share_picture_list_count")
		if err != nil {
			log.Println("GetFileInfoUrl", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			return
		}
		//3. return
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		return
	}
}
