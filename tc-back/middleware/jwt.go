package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tc-back/dfst"
	"time"
)

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "vx:GopherWxf" // 签名信息应该设置成动态从库中获取
)

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}
func GetSignKey() string {
	return SignKey
}
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

type CustomClaims struct {
	Name string `json:"userName"`
	// StandardClaims结构体实现了Claims接口(Valid()函数)
	jwt.StandardClaims
}

//创建一个token
func GenerateToken(c *gin.Context, loginReq dfst.LoginInfoReq) (token string, err error) {
	// 构造SignKey: 签名和解签名需要使用一个值
	jwt2 := NewJWT()
	// 构造用户claims信息(负荷)
	claims := CustomClaims{
		Name: loginReq.User,

		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000, // 签名生效时间
			ExpiresAt: time.Now().Unix() + 3600, // 签名过期时间
			Issuer:    "wxf.top",                // 签名颁发者
		},
	}
	// 根据claims生成token对象
	return jwt2.CreateToken(claims)
}

//创建token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//获取完整的签名令牌
	return token.SignedString(j.SigningKey)
}

//校验token的中间件
func JWTAuth(c *gin.Context) {
	//从header的头部的token字段获取token
	token := c.Request.Header.Get("token")
	log.Println("int jwt  token : ", token)
	if token == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		c.Abort()
		return
	}
	log.Println("recv tokens:", token)
	j := NewJWT()
	//解析token，并将PAYLOAD负载提取出来
	claims, err := j.ParserToken(token)
	if err != nil {
		// token过期
		if err == TokenExpired {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
			//中断调用链
			c.Abort()
			return
		}
		// 其他错误
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
		})
		c.Abort()
		return
	}
	//将负载添加到context上下文中供调用链中的函数使用
	c.Set("claims", claims)
	log.Println("out jwt")
}

//解析token，并将PAYLOAD负载提取出来
func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		// jwt.ValidationError 是一个无效token的错误结构
		if ve, ok := err.(*jwt.ValidationError); ok {
			// ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
				// ValidationErrorExpired表示Token过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
				// ValidationErrorNotValidYet表示无效token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
