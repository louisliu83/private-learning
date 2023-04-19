package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"pa.cn/fedlearn/psi/auth"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/model"
)

// Flags definition
var (
	conf            string
	user            string
	password        string
	party           string
	utype           string
	durationInHours int
)

func init() {
	flag.StringVar(&conf, "c", "conf.json", "set configuration path")
	flag.StringVar(&user, "u", "", "user")
	flag.StringVar(&password, "w", "", "password")
	flag.StringVar(&party, "p", "", "party")
	flag.StringVar(&utype, "t", "agent", "user type['agent','user']")
	flag.IntVar(&durationInHours, "d", 720, "token valid duration in hours")
}

func main() {
	flag.Parse()
	c, err := config.Load(conf)
	if err != nil {
		logrus.Fatalln("Failed to read config!", err)
	}
	// logrus.Infoln(c)
	if err = model.Initdb(c.DB.Path); err != nil {
		logrus.Fatalln(err)
	}
	if user == "" || password == "" || party == "" {
		flag.Usage()
		return
	}
	registerUserIfNotExists(user, password, party)
	generateToken(user, party, utype)
}

func registerUserIfNotExists(name, pass, party string) {
	if err := auth.PasswordSecurityCheck(pass); err != nil {
		logrus.Fatalf("Init admin error:%v\n", err)
	}
	u, err := model.GetUserByUserName(name)
	if err == nil && u != nil {
		u.Party = party
		u.UserPass = auth.HashPassword(pass)
		err = model.UpdateUser(u)
		if err != nil {
			logrus.Warningf("Update user error:%v\n", err)
		}
	} else {
		u := &model.User{
			UserName: name,
			UserPass: auth.HashPassword(pass),
			Party:    party,
		}
		err = model.AddUser(u)
		if err != nil {
			logrus.Fatalf("Init admin error:%v\n", err)
		}
	}
}

func generateToken(name, party, utype string) {
	customClaims := auth.CustomClaims{
		Party:  party,
		UserID: name,
		Admin:  false,
		Type:   utype,
	}
	tokenStr, err := auth.GenerateToken(customClaims, time.Duration(durationInHours)*time.Hour)
	if err != nil {
		logrus.Errorln(err)
	}
	fmt.Printf("Token generated for:%s\n", name)
	fmt.Println(tokenStr)
}
