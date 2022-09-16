package config

import (
	"gopkg.in/ini.v1"
	"log"
	"tc-back/dfst"
)

/* ********************* mysql ************************ */

//解析ini文件并反射到结构体中
func parserDBConfig(dbCfg *dfst.DBConfig) {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		log.Panicf("Fail to read file: %v\n", err)
	}
	dbCfg.Host = cfg.Section("DB").Key("host").String()
	dbCfg.Port, _ = cfg.Section("DB").Key("port").Int()
	dbCfg.User = cfg.Section("DB").Key("user").String()
	dbCfg.Pwd = cfg.Section("DB").Key("pwd").String()
	dbCfg.Database = cfg.Section("DB").Key("database").String()
}

//读取配置文件内容
func LoadDBConfig() *dfst.DBConfig {
	dbCfg := &dfst.DBConfig{}
	//解析ini文件并反射到结构体中
	parserDBConfig(dbCfg)
	return dbCfg
}

/* ********************* redis ************************ */

//解析ini文件并反射到结构体中
func parserRedisConfig(redisCfg *dfst.RedisConfig) {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		log.Panicf("Fail to read file: %v\n", err)
	}
	redisCfg.Host = cfg.Section("Redis").Key("host").String()
	redisCfg.Port, _ = cfg.Section("Redis").Key("port").Int()
	redisCfg.User = cfg.Section("Redis").Key("user").String()
	redisCfg.Pwd = cfg.Section("Redis").Key("pwd").String()
}

//读取配置文件内容
func LoadRedisConfig() *dfst.RedisConfig {
	redisCfg := &dfst.RedisConfig{}
	//解析ini文件并反射到结构体中
	parserRedisConfig(redisCfg)
	return redisCfg
}

/* ********************* front ************************ */
//解析ini文件并反射到结构体中
func parserFrontConfig(frontCfg *dfst.RrontConfig) {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		log.Panicf("Fail to read file: %v\n", err)
	}
	frontCfg.Root = cfg.Section("tc-front").Key("root").String()
	frontCfg.Host = cfg.Section("tc-front").Key("host").String()
	frontCfg.Port = cfg.Section("tc-front").Key("port").String()
}

//读取配置文件内容
func LoadFrontConfig() *dfst.RrontConfig {
	frontCfg := &dfst.RrontConfig{}
	//解析ini文件并反射到结构体中
	parserFrontConfig(frontCfg)
	return frontCfg
}

/* ********************* storage ************************ */
//解析ini文件并反射到结构体中
func LoadStorageConfig() *dfst.StorageConfig {

	cfg, err := ini.ShadowLoad("./config/config.ini")
	if err != nil {
		log.Panicf("Fail to read file: %v\n", err)
	}
	temp := dfst.StorageConfig{GroupToIP: map[string][]string{}}
	temp.GroupToIP["group1"] = cfg.Section("storage").Key("Group1").ValueWithShadows()
	return &temp
}

/* ********************* Fdfs-Cli ************************ */

func LoadFdfsCliConfig() *dfst.FdfsClientConfig {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		log.Panicf("Fail to read file: %v\n", err)
	}
	cli := &dfst.FdfsClientConfig{}
	cli.Conf = cfg.Section("fastdfs-client").Key("conf").String()
	return cli
}
