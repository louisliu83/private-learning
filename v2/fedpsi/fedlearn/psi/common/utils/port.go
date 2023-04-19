package utils

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

const (
	LSOF = "lsof"
)

func CheckPortInUse(port int32) bool {
	cmd := exec.Command(LSOF, fmt.Sprintf("-i:%d", port))
	err := cmd.Run()
	if err != nil {
		logrus.Debugf("Exec command %s failed:%v", cmd.String(), err)
		return false
	}
	return true
}
