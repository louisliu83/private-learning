package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"pa.cn/fedlearn/psi/cmd/controller/bootstrap"
	"pa.cn/fedlearn/psi/config"
)

// Flags definition
var (
	conf  string
	debug bool
)

func init() {
	flag.StringVar(&conf, "c", "conf.json", "set configuration path")
	flag.BoolVar(&debug, "d", false, "set log debug level")
}

func setLogSettings(debug bool) {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	file, err := os.OpenFile("psi.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logrus.Fatalln("Failed to open log file!", err)
	}
	logrus.SetOutput(file)
	logrus.Infoln("##################### PSI Controller Start #####################")
}

func main() {
	// Parse the command-line arguments
	flag.Parse()

	// Set log settings
	setLogSettings(debug)

	// signals
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logrus.Infof("##################### PSI Controller Stop: [%s] #####################", sig)
		done <- true
	}()

	// Load the config
	c, err := config.Load(conf)
	if err != nil {
		logrus.Fatalln("Failed to read config!", err)
	}
	logrus.Debugln(c.String())

	// Boot the platform
	bootstrap.Bootstrap(done)
}
