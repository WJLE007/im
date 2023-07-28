package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"im/define"
	"im/helper"
	"im/models"
	"log"
	"net/http"
	"regexp"
	"time"
)

// Login
//
//	@Description: 登录功能的实现，通过账号和密码登录之后返回token
//	@param c
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

// UserDetail
//
//	@Description: 从token中用户id后获取用户详情
//	@param c
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

// SendCode
//
//	@Description: 发送验证码
//	@param c
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱不能为空",
		})
		return
	}
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$`)
	// 使用正则表达式进行匹配
	match := regex.MatchString(email)
	if !match {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱格式错误",
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
		"msg":  "Success",
	})

	err = models.Redis.Set(context.Background(), define.RegisterPer+email, getRand, time.Second*time.Duration(define.ExpireTime)).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
		})
		log.Println("[DB ERROR]", err)
		return
	}
}

// Register
//
//	@Description:注册接口，对穿回来的参数进行注册服务
//	@param c
func Register(c *gin.Context) {
	code := c.PostForm("code")
	email := c.PostForm("email")
	account := c.PostForm("account")
	password := c.PostForm("password")
	if code == "" || email == "" || account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	//邮箱是否合法（包括格式以及是否存在）
	// 编译正则表达式
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$`)
	// 使用正则表达式进行匹配
	match := regex.MatchString(email)
	if !match {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱格式错误",
		})
		return
	}
	//先对账号和密码进行校验，进行MD5加密
	byAccount, err := models.GetUserBasicCountByAccount(account)
	if err != nil {
		log.Println("[DB ERROR]", err)
		c.JSON(http.StatusOK, gin.H{

			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	if byAccount > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号已经存在",
		})
		return
	}
	//验证码是否正确
	RedisCode, err := models.Redis.Get(context.Background(), define.RegisterPer+email).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	if RedisCode != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}

	ub := &models.UserBasic{
		Identity:  helper.GetUUID(),
		Account:   account,
		Password:  helper.GetMd5(password),
		Email:     email,
		CreatAt:   time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	err = models.InsertOneUserBasic(ub)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误",
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

// UserQuery
//
//	@Description: 获取指定用户的个人信息
//	@param c
func UserQuery(c *gin.Context) {

	account := c.Query("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号不能为空",
		})
	}
	byAccount, err := models.GetUserBasicByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发生系统错误" + err.Error(),
		})
		return
	}
	get := c.MustGet("user_claims").(*helper.UserClaim)
	byAccount.IsFriend, err = models.JudgeIsFriend(get.Identity, byAccount.Identity)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": byAccount,
	})
}

func UserAdd(c *gin.Context) {
	account := c.PostForm("account")
	basicByAccount, err := models.GetUserBasicByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	userClaim := c.MustGet("user_claims").(*helper.UserClaim)
	isFriend, err := models.JudgeIsFriend(basicByAccount.Identity, userClaim.Identity)
	if err != nil {
		log.Println("[DB ERROR]", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库错误",
		})
		return
	}
	if isFriend {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "已经是好友了，不能重复添加",
		})
		return
	}
	//保存房间的记录
	roomBasic := &models.RoomBasic{
		Identity:     helper.GetUUID(),
		UserIdentity: userClaim.Identity,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	if err = models.InsertOneRoomBasic(roomBasic); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库插入错误",
		})
		return
	}
	m := &models.UserRoom{
		UserIdentity: userClaim.Identity,
		RoomIdentity: roomBasic.Identity,
		RoomType:     1,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = models.InsertOneUserRoom(m)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Insert UserRoom Error" + err.Error(),
		})
		return
	}
	m = &models.UserRoom{
		UserIdentity: basicByAccount.Identity,
		RoomIdentity: roomBasic.Identity,
		RoomType:     1,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = models.InsertOneUserRoom(m)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Insert UserRoom Error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "Success",
	})
}
