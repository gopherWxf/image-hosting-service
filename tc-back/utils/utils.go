package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tedcy/fdfs_client"
	"log"
	"strings"
	"tc-back/config"
	"tc-back/dfst"
	"tc-back/opdb"
	"tc-back/oprds"
)

var FrontCfg *dfst.RrontConfig
var StorageCfg *dfst.StorageConfig
var client *fdfs_client.Client
var FdfsCfg *dfst.FdfsClientConfig
var err error

func Default404Router(c *gin.Context) {
	c.File(FrontCfg.Root + "index.html")
	return
}

func LoadConfigAndConn() {
	//从配置文件中读取数据库的配置信息并连接数据库
	if err := opdb.InitMySqlConn(); err != nil {
		log.Panicln(err)
	}
	if err := oprds.InitRedisConn(); err != nil {
		log.Panicln(err)
	}
	//front读取配置文件内容
	FrontCfg = config.LoadFrontConfig()
	//storage
	StorageCfg = config.LoadStorageConfig()
	//fdfs-cli
	client, err = InitFdfsClient()
	if err != nil {
		log.Panicln(err)
	}
}

func InitFdfsClient() (*fdfs_client.Client, error) {
	FdfsCfg = config.LoadFdfsCliConfig()
	temp, err := fdfs_client.NewClientWithConfig(FdfsCfg.Conf)
	return temp, err
}

func UploadFileByBuffer(filename, fileExt string, buf []byte) (string, error) {
	//改成config.ini
	return client.UploadByBuffer(buf, fileExt)
}
func MakeFileURL(fileID string) string {
	//group1/M00/00/00/wKhtZWMgvXuACPLpAAAACa3XZrY899.jpg
	temp := strings.Split(fileID, "/")
	group := temp[0]

	storageIP := StorageCfg.GroupToIP[group][StorageCfg.Cnt]
	StorageCfg.Cnt = (StorageCfg.Cnt + 1) % len(StorageCfg.GroupToIP[group])

	return fmt.Sprintf("http://%s/%s", storageIP, fileID)
}
func DelFileByFileID(fileid string) error {
	if err := client.DeleteFile(fileid); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
