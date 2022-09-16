package opdb

import (
	"errors"
	"tc-back/dfst"
)

/* *********************** Insert *********************** */

func InsertRecordToUserFileList(user, filename, md5 string) error {
	userFileList := dfst.UserFileList{
		User:         user,
		Md5:          md5,
		FileName:     filename,
		SharedStatus: "0",
		Pv:           0,
	}
	return DB.Model(&dfst.UserFileList{}).Create(&userFileList).Error
}

func InsertRecordToFileInfo(md5, fileid, url, fileExt string, size, count int) error {
	fileInfo := dfst.FileInfo{
		Md5:    md5,
		FileId: fileid,
		Url:    url,
		Type:   fileExt,
		Size:   size,
		Count:  count,
	}
	return DB.Model(&dfst.FileInfo{}).Create(&fileInfo).Error
}

func InsertRecordToShareFileList(user, md5, file string) error {
	temp := dfst.ShareFileList{
		User:     user,
		Md5:      md5,
		FileName: file,
		Pv:       0,
	}
	return DB.Model(&dfst.ShareFileList{}).Create(&temp).Error
}

func InsertRecordToSharePicList(user, filemd5, filename, urlmd5 string) error {
	cnt := dfst.SharePictureList{
		User:     user,
		FileMd5:  filemd5,
		FileName: filename,
		UrlMd5:   urlmd5,
		Pv:       0,
	}
	return DB.Model(&dfst.SharePictureList{}).Create(&cnt).Error
}

func InsertRecordToUserInfo(userInfo dfst.RegisterInfoReq) error {
	//查看该用户是否在数据库中
	if HasUserInUserInfo(userInfo.UserName) || HasUserInUserInfo(userInfo.NickName) {
		return errors.New("用户已经存在,请直接登陆")
	}
	user := dfst.UserInfo{
		UserName: userInfo.UserName,
		NickName: userInfo.NickName,
		PassWord: userInfo.FirstPwd,
		Phone:    userInfo.Phone,
		Email:    userInfo.Email,
	}
	//将用户信息插入数据库中
	return DB.Model(&dfst.UserInfo{}).Create(&user).Error
}
