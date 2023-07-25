package service

import (
	"github.com/gin-gonic/gin"
	"im/helper"
	"im/models"
	"log"
	"net/http"
)

func Login(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	if account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或者密码不能为空",
		})
		return
	}
	ub, err := models.GetUserBasicByAP(account, helper.GetMd5(password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或者密码错误",
		})
		return
	}
	token, err := helper.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发生系统错误" + err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": gin.H{
			"token": token,
		},
	})
}

func UserDetail(c *gin.Context) {
	u, _ := c.Get("user_claims")
	uc := u.(*helper.UserClaim)
	userBasic, err := models.GetUserBasicById(uc.Identity)
	if err != nil {
		log.Printf("[DB ERROR]%v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  "Success",
		"data": userBasic,
	})
}

func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱不能为空",
		})
		return
	}
	countByEmail, err := models.GetUserBasicCountByEmail(email)
	if err != nil {
		log.Printf("[DB ERROR]%v\n", err)
		return

	}
	if countByEmail > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "当前邮箱已被注册",
		})
		return
	}
	getRand := helper.GetRand()
	helper.SendCodeEmail(email, getRand)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "seccess",
	})
}
