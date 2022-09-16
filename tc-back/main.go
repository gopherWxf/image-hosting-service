package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ctl "tc-back/controller"
	"tc-back/utils"
)

func main() {
	utils.LoadConfigAndConn()

	r := gin.Default()
	ctl.InitRouter(r)

	r.Run(fmt.Sprintf("%s:%s", utils.FrontCfg.Host, utils.FrontCfg.Port))
}
