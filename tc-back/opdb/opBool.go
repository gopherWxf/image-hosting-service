package opdb

import "tc-back/dfst"

/* *********************** Bool *********************** */

func HasFileInFileInFoByMd5(md5 string) bool {
	var fileInfo dfst.FileInfo
	dbResult := DB.Where(" md5 =  ?", md5).Find(&fileInfo)
	if dbResult.Error != nil {
		return false
	} else {
		return true
	}
}

func HasFileInUserFileList(user, md5, filename string) bool {
	var userFileList dfst.UserFileList
	dbResult := DB.Where("user = ? and md5 =  ? and file_name = ?", user, md5, filename).Find(&userFileList)
	if dbResult.Error != nil {
		return false
	} else {
		return true
	}
}

//查看该用户是否在数据库中
func HasUserInUserInfo(username string) bool {
	result := false
	// 指定库
	var user dfst.UserInfo
	dbResult := DB.Where("user_name =  ?", username).Find(&user)
	if dbResult.Error != nil {
		result = false
	} else {
		result = true
	}
	return result
}

//验证账号密码是否正确，即是否在数据库中存在，注册过
func HasFileInUserInfo(login dfst.LoginInfoReq) (pass bool, err error) {
	var user dfst.UserInfo
	dbErr := DB.Where("user_name = ? && password = ?", login.User, login.Pwd).Find(&user).Error
	if dbErr != nil {
		return false, dbErr
	}
	return true, nil
}
