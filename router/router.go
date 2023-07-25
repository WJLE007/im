package router

import (
	"github.com/gin-gonic/gin"
	"im/middlewares"
	"im/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	//用户登录
	r.POST("/login", service.Login)
	r.POST("/send/code", service.SendCode)

	auth := r.Group("/u", middlewares.AuthCheck())

	auth.GET("/user/detail", service.UserDetail)
	//用户接受和发送消息
	auth.GET("/websocket/message", service.WebsocketMessage)
	return r
}
