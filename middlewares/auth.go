package middlewares

import (
	"github.com/gin-gonic/gin"
	"im/helper"
	"net/http"
)

// AuthCheck
//
//	@Description: 一个中间件，用于鉴别用户的登录状态
//	@return gin.HandlerFunc
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		userClaim, err := helper.AnalyseToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户认证未通过",
			})
			return
		}

		c.Set("user_claims", userClaim)
		c.Next()
	}
}
