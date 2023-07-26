package service

import (
	"github.com/gin-gonic/gin"
	"im/helper"
	"im/models"
	"net/http"
	"strconv"
)

func CharList(c *gin.Context) {
	roomIdentity := c.Query("room_identity")
	if roomIdentity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "房间号不能为空",
		})
		return
	}
	//判断用户是否属于该房间
	userClaim := c.MustGet("user_claims").(*helper.UserClaim)
	_, err := models.GetUserRoomByUserIDRoomID(userClaim.Identity, roomIdentity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "非法访问，用户不属于该房间",
		})
		return
	}
	//选择返回消息的页数
	pageIndex, _ := strconv.ParseInt(c.Query("page_index"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 32)
	skip := (pageIndex - 1) * pageSize
	//查找
	messageListByRoomID, err := models.GetMessageListByRoomID(roomIdentity, &pageSize, &skip)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统异常" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "加载成功",
		"data": gin.H{
			"list": messageListByRoomID,
		},
	})
}
