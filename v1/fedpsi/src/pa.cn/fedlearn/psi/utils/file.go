package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"pa.cn/fedlearn/psi/config"
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

// FileMetaInfo ....
func FileMetaInfo(filePath string) (lineCount int64, size int64, md5 string, err error) {
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
	br := bufio.NewReader(f)
	for {
		_, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return -1, -1, "", err
			}
		}
		lineCount++
	}
	md5 = ""
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
	if _, err := os.Lstat(TaskIntersectPath(taskUID)); err != nil {
		return false
	}
	return true
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
