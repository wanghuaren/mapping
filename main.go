package main

import (
	"lufergo/model"
	"lufergo/model/db"
	"lufergo/service"
	"lufergo/uts"
)

func main() {
	model.InitModel()
	service.InitServeice()

	menuOrderStr := "输入命令:\n"
	menuOrderStr += "gc:手动GC\n"
	menuOrderStr += "exit:退出服务\n"
	menuOrderStr += "按Enter键返回菜单\n"
	menuOrderStr += "a:保存Redis到Mysql\n"
	uts.Log(menuOrderStr)
	uts.Log("本机IP:" + uts.GetLocalIP())

	uts.CommandLine(map[string]func(){"a": db.Cache2DB}, menuOrderStr)
}
