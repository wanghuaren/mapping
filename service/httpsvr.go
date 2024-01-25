package service

import (
	"lufergo/model/netspace"
	"lufergo/model/network"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var regisiter map[string]func(c *gin.Context) = map[string]func(c *gin.Context){}

func initHttp() {
	regisiter["search_to_http"] = search_to_http
	regisiter["search_to_http_api_name"] = search_to_http_api_name
	network.InitHttp(regisiter, websockerSvr)
}

func search_to_http(c *gin.Context) {
	var bt = time.Now().UnixMilli()
	b, _ := c.GetRawData()
	_result := netspace.Search(b)

	_result.Data.Info.Took = payTime(c, bt, "search_to_http")
	c.JSON(http.StatusOK, _result)
}

func search_to_http_api_name(c *gin.Context) {
	var bt = time.Now().UnixMilli()
	b, _ := c.GetRawData()
	_result := netspace.SearchCurrApiName(b)
	payTime(c, bt, "search_to_http_api_name")
	c.JSON(http.StatusOK, _result)
}

func payTime(c *gin.Context, bt int64, funcName string) int {
	Log(funcName+" request ip", c.ClientIP())
	var allT = time.Now().UnixMilli() - bt
	if allT > 1000 {
		Log(funcName+":A", allT)
	} else if allT > 500 {
		Log(funcName+":B", allT)
	} else if allT > 100 {
		Log(funcName+":C", allT)
	} else if allT > 50 {
		Log(funcName+":D", allT)
	}
	return int(allT)
}
