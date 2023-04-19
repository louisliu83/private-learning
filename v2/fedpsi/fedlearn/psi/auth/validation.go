package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wumansgy/goEncrypt/rsa"

	"fedlearn/psi/common/config"
	"fedlearn/psi/model"
)

var privateKey = "MIICXQIBAAKBgQDGkT5+ftmVd/Ve0wxttuvWCk0qPdMnxjhjhJahq3wE/faBo3Oc9UuV6W9+Zj5kt6KZoQfGqO25DqJGE7gm/tVcRknD2WcRjpv1WrsHeeoVIUp1/2lDshiOQXYJNX0no14jYDv+IzRD8p5eG+NDFqwXiKX95fylmgnOj0UoaMKgeQIDAQABAoGAFUgb2pLd3xcsRS15d4jTXe1ct9pId0rXYFMlkc4/TImrkdli2r+vijGqsXFj3oeP9cc8fh483EilO72BTyyg0UKcacPo6pTSKBat392M/BxA7VVUNvkUPUBr9WYp+y6YLm87d915C3ZjaszwVVnjlqhf6rBkTC7UOpRGVjBrvwECQQDaNZotlyTjsdmv4ARkROqdbycfTREqyTaa6PaZIMBxdmP/jCIcvbP9KeYIvkY9Isy2OH8relMCVt13E+SdAgMFAkEA6PTOGxhKDZWqNyz/K1YW70NXZQKsgRHhr4uxu0pGzzs+0N9rJxYEm7ESmEdPSJdEjfXWxaAwErouA7/A7OvJ5QJBALOpvrAa6jyvitTMVdFZDPNjOYsEIUZhNZyGg8PAu7KwD9Www8V2TGP9w3EfeSWNKZA/JDXgGcirTN1me6zqoyECQBKdhV7K6Rf+zrRMBzP6VCjYc8JhnVFPEX7KpfA2dkQXEuT0BYcBDms2kirS//XoCJVjgL8YFt9YO1cXWp5UFTUCQQDJCBdEy5YqSTpWuBhyj+1LO2umw0EUzQtNh4R0/iRsleK/5WWEWr5EjQ4uvxdCxwr6MollijrpZPJBlte398uE"

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

	userPass, err := rsa.RsaDecryptByBase64(userPassB64, privateKey)
	if err != nil {
		return ret, false
	}
	fmt.Printf("rsa解密后:\n%s\n", string(userPass))
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
