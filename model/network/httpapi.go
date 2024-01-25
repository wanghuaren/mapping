package network

import (
	"lufergo/uts"
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func InitHttp(regisiter map[string]func(c *gin.Context), websocketCallback func(*websocket.Conn, []byte) []byte) {
	// if IsDebug {
	gin.SetMode(gin.DebugMode)
	// } else {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	websocketSvr = websocketCallback

	regisiter["socket"] = initWebSocket

	initHttpServer(regisiter)
	// initHttpWeb()
}

// func initHttpWeb() {
// 	httpGin := getNewGin(true)

// 	httpGin.StaticFS("/", gin.Dir("./webroot", false))

// 	httpListenerStr := Conf.String("host")
// 	Log("http web", httpListenerStr)

// 	runGin(httpGin, httpListenerStr, true)
// }

func initHttpServer(regisiter map[string]func(c *gin.Context)) {
	httpGin := getNewGin()
	api := httpGin.Group("api")
	v1 := api.Group("v1")
	for k, v := range regisiter {
		if k == "socket" {
			v1.GET(k, v)
		} else {
			v1.POST(k, v)
		}
	}

	httpListenerStr := Conf.String("host") + ":" + Conf.String("gate_port")
	Log("http server", httpListenerStr)

	runGin(httpGin, httpListenerStr)
}

func getNewGin(isWeb ...bool) *gin.Engine {
	httpGin := gin.New()
	// if len(isWeb) < 1 {
	// 	httpGin.Use(timeoutMiddleware())
	// }
	httpGin.Use(accessControlAllowOriginMiddleware())
	httpGin.SetTrustedProxies([]string{"127.0.0.1"})

	return httpGin
}

func runGin(httpGin *gin.Engine, httpListenerStr string, isWeb ...bool) {
	defer uts.ChkRecover()

	if len(isWeb) > 0 {
		httpListenerStr += ":80"
	}
	go httpGin.Run(httpListenerStr)
}

func timeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(10*time.Millisecond),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}

func timeoutResponse(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, gin.H{
		"code": http.StatusGatewayTimeout,
		"msg":  "请重试",
	})
}

func accessControlAllowOriginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
