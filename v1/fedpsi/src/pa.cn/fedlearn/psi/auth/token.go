package auth

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"pa.cn/fedlearn/psi/config"
)

const (
	TokenIssuer    = "PAB PSI Team"
	TokenTypeAgent = "agent"
	TokenTypeUser  = "user"
)

type CustomClaims struct {
	Party  string `json:"party"`
	UserID string `json:"userID"`
	Admin  bool   `json:"admin"`
	Type   string `json:"utype"`
}

type PSICustomClaims struct {
	jwt.StandardClaims
	CustomClaims
}

func GenerateToken(customClaims CustomClaims, d time.Duration) (tokenString string, err error) {
	claims := &PSICustomClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(d).Unix()),
			Issuer:    TokenIssuer,
		},
		CustomClaims{
			Party:  customClaims.Party,
			Admin:  customClaims.Admin,
			UserID: customClaims.UserID,
			Type:   customClaims.Type,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	keyData, err := ioutil.ReadFile(config.GetConfig().TokenSetting.RsaPrivKeyFile)
	if err != nil {
		logrus.Errorf("%v\n", err)
		return "", err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		logrus.Errorf("%v\n", err)
		return "", err
	}
	tokenString, err = token.SignedString(key)
	return tokenString, nil
}

func ParseToken(tokenSrt string) (*CustomClaims, bool, error) {
	key, err := ioutil.ReadFile(config.GetConfig().TokenSetting.RsaPubKeyFile)
	if err != nil {
		logrus.Errorf("%v\n", err)
		return nil, false, err
	}
	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		logrus.Errorf("%v\n", err)
		return nil, false, err
	}
	token, err := jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return parsedKey, nil
	})
	if err != nil {
		logrus.Errorf("%v\n", err)
		return nil, false, err
	}
	mc := token.Claims.(jwt.MapClaims)
	cc := &CustomClaims{
		Party:  mc["party"].(string),
		Admin:  mc["admin"].(bool),
		UserID: mc["userID"].(string),
		Type:   mc["utype"].(string),
	}
	return cc, token.Valid, nil
}

const (
	levelD = iota
	LevelC
	LevelB
	LevelA
	LevelS
)

const (
	minLength = 8
	maxLength = 16
	minLevel  = LevelA
)

func PasswordSecurityCheck(pwd string) error {
	if len(pwd) < minLength {
		return fmt.Errorf("BAD PASSWORD: The password is shorter than %d characters", minLength)
	}
	if len(pwd) > maxLength {
		return fmt.Errorf("BAD PASSWORD: The password is logner than %d characters", maxLength)
	}

	var level int = levelD
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pwd)
		if match {
			level++
		}
	}

	if level < minLevel {
		return fmt.Errorf("The password does not satisfy the current policy requirements. ")
	}
	return nil
}
