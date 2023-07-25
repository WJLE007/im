package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

type UserClaim struct {
	Identity string `json:"identity"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// GetMd5
// 生成 md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var MyKey = []byte("im-weijiale")

// GenerateToken
// 生成token
//
//	func GenerateToken(indentify, username string) (string, error) {
//		userinfo := &UserClaim{
//			Indentify:      indentify,
//			Username:       username,
//			StandardClaims: jwt.StandardClaims{},
//		}
//		claims := jwt.NewWithClaims(jwt.SigningMethodES256, userinfo)
//		signedString, err := claims.SignedString(MyKey)
//		if err != nil {
//			return "", err
//		} else {
//			return signedString, nil
//		}
//	}
//
// 生成token
func GenerateToken(identity, email string) (string, error) {
	//	objectID, err2 := primitive.ObjectIDFromHex(indentity)
	//if err2 != nil {
	//		return "",err2
	//	}
	userinfo := &UserClaim{
		Identity: identity,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1000000).Unix(), //过期时间
			Issuer:    "im",                                       //签名的发行者
		},
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, userinfo)
	signedString, err := claims.SignedString(MyKey)
	if err != nil {
		return "", err
	} else {
		return signedString, nil
	}

}

// AnalyseToken
//
//	@Description:解析token的操作
//	@param tokenString
//	@return *UserClaim
//	@return error
func AnalyseToken(tokenString string) (*UserClaim, error) {
	userclaim := new(UserClaim)
	claims, err := jwt.ParseWithClaims(tokenString, userclaim, func(token *jwt.Token) (interface{}, error) {
		return MyKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, err

	}
	return userclaim, nil
}

// SendCode
// 发送验证码
func SendCodeEmail(useremail, code string) {
	e := email.NewEmail()
	e.From = "weijiale <wjl13145228831@163.com>"
	e.To = []string{useremail}
	//e.Bcc = []string{"test_bcc@example.com"}
	//e.Cc = []string{"test_cc@example.com"}
	e.Subject = "验证码已发送请注意查收"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("您的验证码是:<b>" + code + "<b>")
	//e.Send("smtp.163.com:465", smtp.PlainAuth("", "wjl13145228831@163.com", "ZMBNIKEPCDPVHVAE", "smtp.163.com"))
	e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "wjl13145228831@163.com", "ZMBNIKEPCDPVHVAE", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}

// GetUUID
// 生成唯一码
func GetUUID() string {
	return uuid.NewV4().String()
}

// GetRand
// 生成验证码

func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}
