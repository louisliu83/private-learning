package useradd

import (
	"flag"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"pa.cn/fedlearn/psi/auth"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/model"
)

var Flags []cli.Flag

func init() {
	Flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "PSI_CONF",
			Name:   "conf",
			Value:  "",
			Usage:  "psi-controller configuration file",
		},

		cli.StringFlag{
			EnvVar: "PSI_USER",
			Name:   "user",
			Value:  "",
			Usage:  "user name to register",
		},

		cli.StringFlag{
			EnvVar: "PSI_PASSWORD",
			Name:   "password",
			Value:  "",
			Usage:  "password for the user",
		},

		cli.StringFlag{
			EnvVar: "PSI_PARTY",
			Name:   "party",
			Value:  "",
			Usage:  "party the user owned to",
		},

		cli.StringFlag{
			EnvVar: "PSI_UTYPE",
			Name:   "utype",
			Value:  "agent",
			Usage:  "user type['agent','user']",
		},

		cli.IntFlag{
			EnvVar: "PSI_TOKEN_VALID_DURATION_HOURS",
			Name:   "duration",
			Value:  720,
			Usage:  "token valid duration in hours",
		},

		cli.BoolFlag{
			EnvVar: "PSI_ISADMIN",
			Name:   "isadmin",
			Usage:  "!!!WARNING:Only the root user should be true",
		},
	}
}

func Action(ctx *cli.Context) {
	conf := ctx.String("conf")
	user := ctx.String("user")
	password := ctx.String("password")
	party := ctx.String("party")
	utype := ctx.String("utype")
	durationInHours := ctx.Int("duration")
	isadmin := ctx.Bool("isadmin")

	c, err := config.Load(conf)
	if err != nil {
		logrus.Fatalln("Failed to read config!", err)
	}

	if err = model.Initdb(c.DB.Path); err != nil {
		logrus.Fatalln(err)
	}

	if user == "" || password == "" || party == "" {
		flag.Usage()
		return
	}

	// Register the user given
	registerUserIfNotExists(user, password, party, isadmin)

	// Generate the token
	// TODO: remove this when generating token from UI is supported.
	generateToken(user, party, utype, durationInHours, isadmin)

}

func Command() cli.Command {
	return cli.Command{
		Name:    "useradd",
		Aliases: []string{"useradd"},
		Usage:   "add a user",
		Flags:   Flags,
		Action:  Action,
	}
}

func registerUserIfNotExists(name, pass, party string, isadmin bool) {
	if err := auth.PasswordSecurityCheck(pass); err != nil {
		logrus.Fatalf("user add %s error:%v\n", name, err)
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
			IsRoot:   isadmin,
		}

		err = model.AddUser(u)
		if err != nil {
			logrus.Fatalf("Init admin error:%v\n", err)
		}
	}
}

func generateToken(name, party, utype string, durationInHours int, isadmin bool) {
	customClaims := auth.CustomClaims{
		Party:  party,
		UserID: name,
		Admin:  isadmin,
		Type:   utype,
	}
	tokenStr, err := auth.GenerateToken(customClaims, time.Duration(durationInHours)*time.Hour)
	if err != nil {
		logrus.Errorln(err)
	}
	logrus.Printf("Token generated for:%s\n", name)
	logrus.Println(tokenStr)
}
