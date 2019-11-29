package lib

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
)

var SignKey = "webblog"

//GetUniqueID 返回唯一哈希值
func GetUniqueID() string {
	return xid.New().String()
}

// CheckEmail check if the email is valid
func CheckEmail(email string) bool {
	reg, err := regexp.Compile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if err != nil {
		log.Println(err)
		return false
	}
	ret := reg.Find([]byte(email))
	if ret == nil || string(ret) != email {
		return false
	}
	return true
}

// CheckUsername check username
func CheckUsername(username string) bool {
	reg, err := regexp.Compile(`[a-z|_]{1}[a-z0-9_\-]{4,15}`)
	if err != nil {
		log.Println(err)
		return false
	}
	ret := reg.Find([]byte(username))
	if ret == nil || string(ret) != username {
		return false
	}
	return true
}

//CheckToken 检查用户token是否有效
func CheckToken(token string) (bool, error) {
	ParsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("该token使用算法非HS加密，为：%s", token.Header["alg"])
		}
		return []byte(SignKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	var exp int64
	if claims, ok := ParsedToken.Claims.(jwt.MapClaims); ok && ParsedToken.Valid {
		expired := claims["exp"]
		if expired == nil {
			return false, nil
		}
		exp = int64(expired.(float64))
		if time.Now().Unix() > exp {
			return false, nil
		}
	} else {
		return false, nil
	}
	return true, nil
}

//GenerateToken 为用户生成Token
func GenerateToken(useremail string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), //过期时间 一个月
		"sub": useremail,                                  //主题
		"iat": time.Now().Unix(),                          //发行时间
		"jti": uuid.NewV1(),                               //ID
	})
	return token.SignedString([]byte(SignKey)) //加签名
}

// GetUserEmailFromToken get useremail 'sub' from token
func GetUserEmailFromToken(jwtString string, JWTKey string) (bool, string) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("该token使用算法非HS加密，为：%s", token.Header["alg"])
		}
		return []byte(JWTKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		useremail := claims["sub"]
		if useremail == "" {
			return false, ""
		}
		sub := string(useremail.(string))
		return true, sub
	}
	return false, ""
}
