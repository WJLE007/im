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
	//  用户发送验证码
	r.POST("/send/code", service.SendCode)
	//用户注册
	r.POST("/register", service.Register)
	auth := r.Group("/u", middlewares.AuthCheck())
	//查询用户详情
	auth.GET("/user/detail", service.UserDetail)
	//用户接受和发送消息
	auth.GET("/websocket/message", service.WebsocketMessage)
	//查询消息记录
	auth.GET("/chat/list", service.CharList)
	//查询指定用户的个人信息

	return r
}
