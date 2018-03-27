package main

import (
	"LeafProject/leaf"
	lconf "LeafProject/leaf/conf"
	"server/conf"
	"server/game"
	"server/gate"
	"server/login"
	"LeafProject/leaf/db/mongodb"
	"fmt"
)

func main() {
	//初始化 全局配置
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	//初始化数据库
	c, err := mongodb.Dial(conf.Server.MongodbUrl, conf.Server.MongodbSessionNum)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	//运行模块
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
