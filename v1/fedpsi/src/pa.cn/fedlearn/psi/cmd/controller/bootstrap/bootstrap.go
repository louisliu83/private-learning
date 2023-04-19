package bootstrap

import (
	"os"
	"path/filepath"

	"pa.cn/fedlearn/psi/server/http"
	"pa.cn/fedlearn/psi/server/scheduler"

	"github.com/sirupsen/logrus"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/model"
	"pa.cn/fedlearn/psi/tproxy"
)

// Bootstrap boots system
func Bootstrap(done <-chan bool) {
	initDirs()
	initDB()
	bootTProxy()
	bootHttpServer()
	bootTaskScheduler()
	bootTaskRerunScheduler()
	bootJobScheduler()
	bootDataScheduler()
	bootDatasetSharder()
	bootDatasetDownloader()

	// waiting done signal
	<-done
}

func initDirs() {
	// tasks/jobs directory
	logrus.Infoln("Initialize dirs of tasks/jobs ...")
	if err := os.MkdirAll(config.GetConfig().TaskSetting.TasksDir, os.ModePerm); err != nil {
		logrus.Errorln("Initialize tasks dir error:", err)
		os.Exit(-1)
	}

	// dataset base dirs
	logrus.Infoln("Initialize dirs of dataset ...")
	if err := os.MkdirAll(config.GetConfig().DataSet.Dir, os.ModePerm); err != nil {
		logrus.Errorln("Initialize Dataset dir error:", err)
		os.Exit(-1)
	}

	// dataset base tmp dirs
	logrus.Infoln("Initialize tmp dirs of dataset ...")
	if err := os.MkdirAll(filepath.Join(config.GetConfig().DataSet.Dir, "tmp"), os.ModePerm); err != nil {
		logrus.Errorln("Initialize Dataset tmp dir error:", err)
		os.Exit(-1)
	}

}

func initDB() {
	logrus.Infoln("Initialize db ...")
	if err := model.Initdb(config.GetConfig().DB.Path); err != nil {
		logrus.Errorln("Initialize Database error:", err)
		os.Exit(-1)
	}
}

func bootHttpServer() {
	http.GetServer().Start()
}

func bootTProxy() {
	if config.GetConfig().FeatureGate.TProxy {
		go tproxy.GetTProxy().Start()
	}
}

// Schedulers in background

func bootDataScheduler() {
	scheduler.GetDatasetScheduler().Start()
}

func bootDatasetSharder() {
	scheduler.GetSharderScheduler().Start()
}

func bootDatasetDownloader() {
	scheduler.GetDownloadScheduler().Start()
}

func bootJobScheduler() {
	scheduler.GetJobScheduler().Start()
}

func bootTaskScheduler() {
	scheduler.GetTaskScheduler().Start()
}

func bootTaskRerunScheduler() {
	if config.GetConfig().TaskSetting.AutoRerun {
		scheduler.GetTaskRerunScheduler().Start()
	}
}
