package auth

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/model"
)

// AuthenticatedInfo is the info of authenticated user
type AuthenticatedInfo struct {
	UserName string
	Party    string
	IsAdmin  bool
	Type     string
}

type authenticateFunc func(string) (AuthenticatedInfo, bool)

var authenticateFuncs = map[string]authenticateFunc{}

func init() {
	authenticateFuncs["Basic"] = basicAuthenticate
	authenticateFuncs["Bearer"] = bearerAuthenticate
}

func Authenticate(credential string) (AuthenticatedInfo, bool) {
	ret := AuthenticatedInfo{}
	credArray := strings.Split(credential, " ")
	if len(credArray) != 2 {
		return ret, false
	}
	authType := credArray[0]
	c := credArray[1]
	if f, ok := authenticateFuncs[authType]; ok {
		return f(c)
	}
	return ret, false
}

func basicAuthenticate(userPassB64 string) (AuthenticatedInfo, bool) {
	ret := AuthenticatedInfo{}
	userPass, err := base64.StdEncoding.DecodeString(userPassB64)
	if err != nil {
		return ret, false
	}
	userPwd := strings.SplitN(string(userPass), ":", 2)
	if len(userPwd) != 2 {
		return ret, false
	}
	username := userPwd[0]
	u, err := model.GetUserByUserName(username)
	if err != nil {
		return ret, false
	}

	presentDate := time.Now().Format("2006-01-02")
	if presentDate != u.FailedDate {
		// reset next day
		u.FailedDate = presentDate
		u.FailedCount = 0
		if err = model.UpdateUser(u); err != nil {
			logrus.Errorf("Update user login failed count err:%v", err)
		}
	}

	if u.FailedCount >= model.LockIfFailedCountPerDay {
		// return if failed count exceed LockCount
		return ret, false
	}

	password := userPwd[1]
	hashedPass := HashPassword(password)
	if u.UserPass != hashedPass {
		// add failed count if failed
		u.FailedCount++
		if err = model.UpdateUser(u); err != nil {
			logrus.Errorf("Update user login failed count err:%v", err)
		}
		return ret, false
	}
	// reset failed count if true
	u.FailedCount = 0
	if err = model.UpdateUser(u); err != nil {
		logrus.Errorf("Update user login failed count err:%v", err)
	}
	ret.UserName = u.UserName
	ret.Party = u.Party
	if ret.Party == "" {
		//local user
		ret.Party = config.GetConfig().PartyName
	}
	ret.IsAdmin = u.IsRoot
	ret.Type = TokenTypeUser
	return ret, true
}

func bearerAuthenticate(token string) (AuthenticatedInfo, bool) {
	ret := AuthenticatedInfo{}
	cc, valid, err := ParseToken(token)
	if err != nil {
		return ret, false
	}
	if !valid {
		return ret, false
	}
	if cc == nil {
		return ret, false
	}
	ret.UserName = cc.UserID
	ret.Party = cc.Party
	ret.IsAdmin = cc.Admin
	ret.Type = cc.Type
	return ret, true
}
