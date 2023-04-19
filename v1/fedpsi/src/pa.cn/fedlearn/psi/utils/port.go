package utils

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func CheckPortInUse(port int32) bool {
	cmd := exec.Command("lsof", fmt.Sprintf("-i:%d", port))
	err := cmd.Run()
	if err != nil {
		logrus.Infof("Exec cmd %s %v failed:%v", cmd.Path, cmd.Args, err)
		return false
	}
	return true
}
