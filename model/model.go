package model

import (
	"lufergo/model/db"
	"lufergo/model/netspace"
	"lufergo/uts"
)

func InitModel() {
	uts.InitLog("lufer")
	db.InitRedis()
	db.DB2Cache()
	// db.Cache2DB()

	netspace.InitNetSpace()
}
