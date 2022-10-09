package dfst

import "time"

//将结构体抽离出来，结果包循环引用的问题

/* *********************** config struct *********************** */

//mysql config
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
}

//redis config
type RedisConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

//tc-front config
type RrontConfig struct {
	Root string `json:"root"`
	Host string `json:"host"`
	Port string `json:"port"`
}

//fastdfs-client config
type FdfsClientConfig struct {
	Conf string `json:"conf"`
}

//storage config
type StorageConfig struct {
	GroupToIP map[string][]string `json:"grouptoip"`
	Cnt       int                 //轮询，负载均衡
}

/* *********************** gorm struct *********************** */

//user info 用户信息表
type UserInfo struct {
	Id         int32     `gorm:"AUTO_INCREMENT,primary_key"`
	UserName   string    `gorm:"column:user_name"`
	NickName   string    `gorm:"column:nick_name"`
	PassWord   string    `gorm:"column:password"`
	Phone      string    `gorm:"column:phone"`
	Email      string    `gorm:"column:email"`
	CreateTime time.Time `gorm:"-"`
}

//file_info 文件信息表
type FileInfo struct {
	Id     int32  `gorm:"AUTO_INCREMENT,primary_key"`
	Md5    string `gorm:"column:md5"`
	FileId string `gorm:"column:file_id"`
	Url    string `gorm:"column:url"`
	Size   int    `gorm:"column:size"`
	Type   string `gorm:"column:type"`
	Count  int    `gorm:"column:count"`
}

//user_file_list 用户文件列表
type UserFileList struct {
	Id           int32     `gorm:"AUTO_INCREMENT,primary_key"`
	User         string    `gorm:"column:user"`
	Md5          string    `gorm:"column:md5"`
	CreateTime   time.Time `gorm:"-"`
	FileName     string    `gorm:"column:file_name"`
	SharedStatus string    `gorm:"column:shared_status"`
	Pv           int       `gorm:"column:pv"`
}

//user_file_count 用户文件数量表
type UserFileCount struct {
	Id    int32  `gorm:"AUTO_INCREMENT,primary_key"`
	User  string `gorm:"column:user"`
	Count int    `gorm:"column:count"`
}

//share_file_list 共享文件列表
type ShareFileList struct {
	Id         int32     `gorm:"AUTO_INCREMENT,primary_key"`
	User       string    `gorm:"column:user"`
	Md5        string    `gorm:"column:md5"`
	FileName   string    `gorm:"column:file_name"`
	Pv         int       `gorm:"column:pv"`
	CreateTime time.Time `gorm:"-"`
}

//share_picture_list 共享图片列表
type SharePictureList struct {
	Id         int32  `gorm:"AUTO_INCREMENT,primary_key"`
	User       string `gorm:"column:user"`
	FileMd5    string `gorm:"column:filemd5"`
	FileName   string `gorm:"column:file_name"`
	UrlMd5     string `gorm:"column:urlmd5"`
	Pv         int    `gorm:"column:pv"`
	CreateTime string `gorm:"-"`
}

/* *********************** request struct *********************** */

type RegisterInfoReq struct {
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	FirstPwd string `json:"firstPwd"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type LoginInfoReq struct {
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

type MD5Req struct {
	Token    string `json:"token"`
	Md5      string `json:"md5"`
	Filename string `json:"filename"`
	User     string `json:"user"`
}

type MyFilesReq struct {
	Token string `json:"token"`
	User  string `json:"user"`
	Count int    `json:"count"`
	Start int    `json:"start"`
}

type ShareFilesReq struct {
	Count int `json:"count"`
	Start int `json:"start"`
}

type DealFileReq struct {
	Token    string `json:"token"`
	User     string `json:"user"`
	Md5      string `json:"md5"`
	Filename string `json:"filename"`
}

type DealShareFileReq struct {
	User     string `json:"user"`
	Md5      string `json:"md5"`
	Filename string `json:"filename"`
}

type SharePicReq struct {
	Token    string `json:"token"`
	User     string `json:"user"`
	Count    int    `json:"count"`
	Start    int    `json:"start"`
	Md5      string `json:"md5"`
	UrlMd5   string `json:"urlmd5"`
	Filename string `json:"filename"`
}

/* *********************** response struct *********************** */

//bug 前端json是share_status ,刷新不会请求cmd=count
type UserFileListAndFileInfoResp struct {
	User         string    `json:"user"`
	Md5          string    `json:"md5"`
	CreateTime   time.Time `json:"create_time"`
	FileName     string    `json:"file_name"`
	SharedStatus int       `json:"share_status" gorm:"column:shared_status"`
	Pv           int       `json:"pv" gorm:"column:pv"`

	Url  string `json:"url"`
	Size string `json:"size"`
	Type string `json:"type"`
}

type ShareFileListAndFileInfoResp struct {
	User         string    `json:"user"`
	Md5          string    `json:"md5"`
	FileName     string    `json:"file_name"`
	CreateTime   time.Time `json:"create_time"`
	SharedStatus int       `json:"share_status" gorm:"column:shared_status"`
	Pv           int       `json:"pv" gorm:"column:pv"`

	Url  string `json:"url"`
	Size string `json:"size"`
	Type string `json:"type"`
}

type ShareFileListResp struct {
	Md5      string `json:"-"`
	FileName string `json:"file_name"`
	Pv       int    `json:"pv"`
}

//buf 下载榜是filename
type ShareFileList2Resp struct {
	FileName string `json:"filename"`
	Pv       int    `json:"pv"`
}

type SharePicListAndFileInfoResp struct {
	User       string    `json:"user" gorm:"column:user"`
	FileMd5    string    `json:"filemd5" gorm:"column:filemd5"`
	UrlMd5     string    `json:"urlmd5" gorm:"column:urlmd5"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	FileName   string    `json:"file_name" gorm:"column:file_name"`
	Pv         int       `json:"pv" gorm:"column:pv"`
	Size       string    `json:"size" gorm:"column:size"`
}

//share_picture_list 共享图片列表
type SharePicListResp struct {
	Code       int    `json:"code"`
	User       string `json:"user"`
	Pv         int    `json:"pv"`
	CreateTime string `json:"time"`
	Url        string `json:"url"`
}
