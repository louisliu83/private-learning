package bootstrap

import (
	"os"
	"path/filepath"

	"fedlearn/psi/server/http"
	"fedlearn/psi/server/scheduler"

	"fedlearn/psi/common/config"
	"fedlearn/psi/common/portmgr"
	"fedlearn/psi/model"
	grpcproxy "fedlearn/psi/proxy/grpc"

	"github.com/sirupsen/logrus"
)

// Bootstrap boots system
func Bootstrap(done <-chan bool) {
	initDirs()
	initDB()
	initPortManager()
	bootControllerServer()
	bootDataServer()
	bootTaskScheduler()
	bootTaskRerunScheduler()
	bootJobScheduler()
	bootActivityScheduler()
	bootDataScheduler()
	bootDatasetSharder()
	bootDatasetDownloader()
	bootTaskStatusScheduler()

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

func initPortManager() {
	portmgr.PSIServerPortManager = portmgr.NewPortManager(37700, config.GetConfig().PsiExecutor.ServerConcurrency*2)
	portmgr.PSIClientPortManager = portmgr.NewPortManager(27700, config.GetConfig().PsiExecutor.ClientConcurrency*2)
}

func initDB() {
	logrus.Infoln("Initialize db ...")
	if err := model.Initdb(config.GetConfig().DB.Path); err != nil {
		logrus.Errorln("Initialize Database error:", err)
		os.Exit(-1)
	}
}

func bootControllerServer() {
	http.GetServer().Start()
}

func bootDataServer() {
	go func() {
		server := grpcproxy.NewGrpc2TCPServer(config.GetConfig().StreamListener.Address)
		server.Start()
	}()
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

func bootActivityScheduler() {
	scheduler.GetActivityScheduler().Start()
}

func bootTaskScheduler() {
	scheduler.GetTaskScheduler().Start()
}

func bootTaskRerunScheduler() {
	if config.GetConfig().TaskSetting.AutoRerun {
		scheduler.GetTaskRerunScheduler().Start()
	}
}

func bootTaskStatusScheduler() {
	scheduler.GetTaskStatusScheduler().Start()
}
