package opdb

import (
	"errors"
	"log"
	"tc-back/dfst"
)

/* *********************** Select *********************** */

func GetSharePicListJoinFileInfo(user string, start, count int) ([]dfst.SharePicListAndFileInfoResp, error) {
	var sharePicInfos []dfst.SharePicListAndFileInfoResp

	DB.Raw("select "+
		"share_picture_list.user, "+
		"share_picture_list.filemd5, "+
		"share_picture_list.file_name,"+
		"share_picture_list.urlmd5, "+
		"share_picture_list.pv, "+
		"share_picture_list.create_time, "+
		"file_info.size "+
		"from file_info, share_picture_list "+
		"where share_picture_list.user = ? and  file_info.md5 = share_picture_list.filemd5 limit ?, ? ",
		user, start, count).Scan(&sharePicInfos)
	if sharePicInfos == nil {
		return nil, errors.New("GetSharePicListJoinFileInfo : not found")
	}
	return sharePicInfos, nil
}
func GetFileInfoCount(md5 string) (int, error) {
	var fileInfo dfst.FileInfo
	dbResult := DB.Where(" md5 =  ?", md5).Find(&fileInfo).Error
	return fileInfo.Count, dbResult
}
func GetFileInfoFileId(md5 string) (string, error) {
	var info dfst.FileInfo
	err := DB.Where("md5 = ?", md5).Find(&info).Error
	if err != nil {
		return "", err
	}
	return info.FileId, nil
}
func GetSharePicListInfo(urlmd5 string) (dfst.SharePictureList, error) {
	var info dfst.SharePictureList
	err := DB.Where("urlmd5 = ?", urlmd5).Find(&info).Error
	if err != nil {
		return info, err
	}
	return info, nil
}

func GetFileInfoUrl(md5 string) (string, error) {
	var fileInfo dfst.FileInfo
	dbResult := DB.Where(" md5 =  ?", md5).Find(&fileInfo).Error
	return fileInfo.Url, dbResult
}

func GetUserFileCount(user string) (int, error) {
	var userfilecount dfst.UserFileCount
	dbResult := DB.Where(" user =  ?", user).Find(&userfilecount)
	if dbResult.RowsAffected == 0 {
		return 0, nil
	} else if dbResult.Error != nil {
		log.Println("GetUserFileCount err", dbResult.Error)
		return 0, dbResult.Error
	} else {
		return userfilecount.Count, nil
	}
}
func GetUserFileListJoinFileInfo(user, cmd string, start, count int) ([]dfst.UserFileListAndFileInfoResp, error) {
	var fileInfos []dfst.UserFileListAndFileInfoResp
	if cmd == "normal" {
		DB.Raw("select user_file_list.*, file_info.url, file_info.size, file_info.type from file_info, user_file_list where user = ? and file_info.md5 = user_file_list.md5 limit ?, ? ", user, start, count).Scan(&fileInfos)
	} else if cmd == "pvasc" {
		DB.Raw("select user_file_list.*, file_info.url, file_info.size, file_info.type from file_info, user_file_list where user = ? and file_info.md5 = user_file_list.md5 order by pv asc  limit ?, ? ", user, start, count).Scan(&fileInfos)
	} else {
		DB.Raw("select user_file_list.*, file_info.url, file_info.size, file_info.type from file_info, user_file_list where user = ? and file_info.md5 = user_file_list.md5 order by pv desc  limit ?, ? ", user, start, count).Scan(&fileInfos)
	}
	if fileInfos == nil {
		return nil, errors.New("GetUserFileListJoinFileInfo : not found")
	}
	return fileInfos, nil
}

func GetShareFileListJoinFileInfo(start, count int) ([]dfst.ShareFileListAndFileInfoResp, error) {
	var sharefileInfos []dfst.ShareFileListAndFileInfoResp
	DB.Raw("select share_file_list.*, file_info.url, file_info.size, file_info.type from file_info, share_file_list where file_info.md5 = share_file_list.md5 limit ?, ? ", start, count).Scan(&sharefileInfos)
	if sharefileInfos == nil {
		return nil, errors.New("GetShareFileListJoinFileInfo : not found")
	}
	return sharefileInfos, nil
}

func GetShareFileList() ([]dfst.ShareFileListResp, error) {
	var sharefile []dfst.ShareFileListResp

	DB.Raw("select md5, file_name, pv from share_file_list order by pv desc").Scan(&sharefile)
	if sharefile == nil {
		return nil, errors.New("GetShareFileList : not found")
	}
	return sharefile, nil
}
