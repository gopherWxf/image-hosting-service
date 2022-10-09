package opdb

import (
	"github.com/jinzhu/gorm"
	"tc-back/dfst"
)

/* *********************** Update *********************** */

func UpdateUserFileListForShare(user, md5, file string) error {
	return DB.Exec("update user_file_list set shared_status = 1 where user = ? and md5 = ? and file_name = ?", user, md5, file).Error
}
func UpdateUserFileCountDecr(user string) error {
	return DB.Model(&dfst.UserFileCount{}).Where("user = ?", user).Update("count = count-1").Error
}

func UpdateFileInfoCountDecr(md5 string) error {
	return DB.Exec("update file_info set count= count - 1 where md5 = ?", md5).Error
}

func UpdateFileInfoCountInc(md5 string) error {
	return DB.Exec("update file_info set count= count + 1 where md5 = ?", md5).Error
}

func UpdateUserFileListPvInc(user, md5, filename string) {
	DB.Exec("update user_file_list set pv = pv + 1 where user = ? and md5 = ? and file_name = ? ", user, md5, filename)
}

func UpdateUserFileListForNotShare(user, md5, file string) error {
	return DB.Exec("update user_file_list set shared_status = 0 where user = ? and md5 = ? and file_name = ?", user, md5, file).Error
}

func UpdateShareFileListPvInc(md5, filename string) {
	DB.Exec("update share_file_list set pv = pv + 1 where md5 = ? and file_name = ? ", md5, filename)
}
func UpdateSharePicListPvInc(md5 string) {
	DB.Exec("update share_picture_list set pv = pv + 1 where urlmd5 = ?  ", md5)
}
func UpdateUserFileListCountInc(user string) error {
	var userFileListCount dfst.UserFileCount
	dbResult := DB.Where(" user =  ?", user).Find(&userFileListCount)
	//不存在则创建
	if dbResult.RowsAffected == 0 {
		userFileListCount.User = user
		userFileListCount.Count = 1
		return DB.Model(&dfst.UserFileCount{}).Create(&userFileListCount).Error

	} else {
		return DB.Model(&dfst.UserFileCount{}).Where("user = ?", user).Update("count", gorm.Expr("count+1")).Error
	}
}
