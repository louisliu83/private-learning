package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	TimeFormatString = "2006-01-02 15:04:05"
)

var (
	initialized bool
	c           Config
)

type Listener struct {
	Address     string `json:"address"`
	TLSEnabled  bool   `json:"tlsEnabled"`
	TLSAddress  string `json:"tlsAddress"`
	TLSCertFile string `json:"tlsCertFile"`
	TLSKeyFile  string `json:"tlsKeyFile"`
}

type StreamListener struct {
	Address string `json:"address"`
}

type DB struct {
	Path string `json:"path"`
}

type DataSet struct {
	ValidDays   int32  `json:"validDays"`
	MaxLines    int64  `json:"maxLine"`
	Dir         string `json:"dir"`
	Sharders    int    `json:"sharders"`
	Downloaders int    `json:"downloaders"`
}

type TaskSetting struct {
	TasksDir  string `json:"dir"`
	AutoRerun bool   `json:"autoRerun"`
}

type PartyInfo struct {
	Name    string `json:"name"`
	Scheme  string `json:"scheme"`
	Address string `json:"address"`
}

type PsiExecutor struct {
	BinPath           string `json:"binpath"`
	ServerTimeout     int32  `json:"serverTimeout"`
	ServerConcurrency int    `json:"serverConcurrency"`
	ClientConcurrency int    `json:"clientConcurrency"`
}

type TokenSetting struct {
	AuthEnabled       bool   `json:"authEnabled"`
	TokenValidInHours int32  `json:"tokenValidInHours"`
	RsaPubKeyFile     string `json:"rsaPubKeyFile"`
	RsaPrivKeyFile    string `json:"rsaPrivKeyFile"`
}

type FeatureGate struct {
	DatasetPrivate bool `json:"datasetPrivate"`
	Audit          bool `json:"audit"`
}

type Audit struct {
	File string `json:"file"`
}

type Config struct {
	PartyName      string         `json:"partyname"`
	FeatureGate    FeatureGate    `json:"featureGate"`
	Listener       Listener       `json:"listen"`
	StreamListener StreamListener `json:"streamListen"`
	DB             DB             `json:"db"`
	DataSet        DataSet        `json:"dataset"`
	TaskSetting    TaskSetting    `json:"tasks"`
	PsiExecutor    PsiExecutor    `json:"psiExecutor"`
	TokenSetting   TokenSetting   `json:"tokenSetting"`
	Audit          Audit          `json:"audit"`
}

func (c Config) String() string {
	data, _ := json.MarshalIndent(c, "", "\t")
	return string(data)
}

func Load(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	c.setDefaultConfigIfNotSet()
	initialized = true
	return &c, nil
}

func GetConfig() Config {
	if !initialized {
		logrus.Fatalln("Cannot use config before it is initialized.")
		os.Exit(-1)
	}
	return c
}

func GetDatasetValidDays() int32 {
	if GetConfig().DataSet.ValidDays == 0 {
		return defaultDatasetValidDays
	}
	return GetConfig().DataSet.ValidDays
}
