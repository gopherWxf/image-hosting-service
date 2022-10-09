package opdb

/* *********************** Del *********************** */

func DelRecordForFileInfo(md5 string) {
	DB.Exec("delete from file_info where md5 = ?", md5)
}

func DelRecordForShareFileList(user, md5, filename string) error {
	return DB.Exec("delete from share_file_list where user = ? and md5 = ? and file_name = ?", user, md5, filename).Error
}

func DelRecordForSharePicList(user, urlmd5 string) error {
	return DB.Exec("delete from share_picture_list "+
		"where user = ? and urlmd5 = ? ",
		user, urlmd5).Error
}

func DelRecordForUserFileList(user, md5, filename string) error {
	return DB.Exec("delete from user_file_list where user = ? and md5 = ? and file_name = ?", user, md5, filename).Error
}
