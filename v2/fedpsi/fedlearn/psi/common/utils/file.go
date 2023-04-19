package utils

import (
	"bufio"
	"fedlearn/psi/common/config"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	//FileTaskPID : <DIR>/<TaskUUID>/pid Store the task PID
	FileTaskPID = "pid"
	//FileTaskStatus : <DIR>/<TaskUUID>/info Store the task status
	FileTaskStatus = "info"
	//FileTaskResult : <DIR>/<TaskUUID>/result Store the task result
	FileTaskResult = "result"
	//FileTaskIntersect : <DIR>/<TaskUUID>/intersect Store the task result intersect(only in server task)
	FileTaskIntersect = "intersect"
	//FileJobIntersect : <DIR>/<JobUUID>/intersect Store the task result intersect(only in server task)
	FileJobIntersect = "intersect"
)

/* ############################## Dataset file processing ############################## */

// TaskDataSetPath ...
func TaskDataSetPath(md5 string, name string, index int32) string {
	return filepath.Join(config.GetConfig().DataSet.Dir, md5, fmt.Sprintf("%s_%d", name, index))
}

// FileMetaInfo return the lineCount and size of the file
func FileMetaInfo(filePath string) (lineCount, size int64, repeatedLine string, err error) {
	fi, err := os.Lstat(filePath)
	if err != nil {
		return -1, -1, "", err
	}
	size = fi.Size()
	f, err := os.Open(filePath)
	if err != nil {
		return -1, -1, "", err
	}
	defer f.Close()
	lineCount = 0
	counts := make(map[string]int)
	br := bufio.NewReader(f)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return -1, -1, "", err
			}
		}
		lineCount++
		counts[string(line)]++
	}

	for line, n := range counts {
		if n > 1 {
			repeatedLine = line
			return
		}
	}
	return
}

/* ############################## task file processing ############################## */

// TaskPath ...
func TaskPath(taskUID string) string {
	return filepath.Join(config.GetConfig().TaskSetting.TasksDir, taskUID)
}

// TaskPIDPath ...
func TaskPIDPath(taskUID string) string {
	return filepath.Join(TaskPath(taskUID), FileTaskPID)
}

// TaskStatusPath ...
func TaskStatusPath(taskUID string) string {
	return filepath.Join(TaskPath(taskUID), FileTaskStatus)
}

// TaskResultPath ...
func TaskResultPath(taskUID string) string {
	return filepath.Join(TaskPath(taskUID), FileTaskResult)
}

// TaskIntersectPath ...
func TaskIntersectPath(taskUID string) string {
	return filepath.Join(TaskPath(taskUID), FileTaskIntersect)
}

// IsTaskIntersectExists ...
func IsTaskIntersectExists(taskUID string) bool {
	return FileExists(TaskIntersectPath(taskUID))
}

/* ############################## job file processing ############################## */

// JobPath ...
func JobPath(jobUID string) string {
	return filepath.Join(config.GetConfig().TaskSetting.TasksDir, fmt.Sprintf("j_%s", jobUID))
}

// JobIntersectPath ...
func JobIntersectPath(jobUID string) string {
	return filepath.Join(JobPath(jobUID), FileJobIntersect)
}

// IsJobIntersectExists ...
func IsJobIntersectExists(jobUID string) bool {
	if _, err := os.Lstat(JobIntersectPath(jobUID)); err != nil {
		return false
	}
	return true
}

// FileExists check where the filename exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err == nil && !info.IsDir() {
		return true
	} else {
		return false
	}
}

const (
	DOS2UNIX = "dos2unix"
	SORT     = "sort"
)

// DeDuplicate deduplicate the file named fileName
// sort -u -o outputFileName inputFileName
func DeDuplicate(filename string) error {
	if !FileExists(filename) {
		return fmt.Errorf("file %s doesnot exist", filename)
	}
	tmpFileName := fmt.Sprintf("%s.tmp", filename)
	if err := os.Rename(filename, tmpFileName); err != nil {
		return err
	}
	cmd := exec.Command(SORT, "-u", "-o", filename, tmpFileName)
	// var stdout, stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		logrus.Debugf("Exec command %s failed:%v", cmd.String(), err)
		return err
	}
	return nil
}

func Dos2Unix(filename string) error {
	if !FileExists(filename) {
		return fmt.Errorf("file %s doesnot exist", filename)
	}
	cmd := exec.Command(DOS2UNIX, filename)
	// var stdout, stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		logrus.Debugf("Exec command %s failed:%v", cmd.String(), err)
		return err
	}
	return nil
}
