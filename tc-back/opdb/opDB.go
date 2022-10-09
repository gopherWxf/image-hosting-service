package opdb

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"tc-back/config"
)

var DB *gorm.DB

//从配置文件中读取数据库的配置信息并连接数据库
func InitMySqlConn() (err error) {
	//读取配置文件内容
	dbCfg := config.LoadDBConfig()
	//拼凑连接数据库的语句
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbCfg.User, dbCfg.Pwd, dbCfg.Host, dbCfg.Port, dbCfg.Database)
	log.Println(connStr)
	//连接数据库
	DB, err = gorm.Open("mysql", connStr)
	DB.SingularTable(true)

	if err != nil {
		return err
	}
	//检测数据库是否活跃
	return DB.DB().Ping()
}
