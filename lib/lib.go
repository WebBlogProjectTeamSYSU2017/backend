package lib

import (
	"log"
	"regexp"

	"github.com/rs/xid"
)

//GetCheckSumIEEE 返回给定标题的唯一哈希值
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
