package worker

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/model"
	"pa.cn/fedlearn/psi/utils"
)

// MergeJobIntersect merge all tasks' intersect together
// At last, if the update the job's intersect count in db in client side
// TODO, how to sync the intersect file and job intersect count to the server side
func MergeJobIntersect(ctx context.Context, tasks []*model.Task, job *model.Job) error {
	// only merge job intersect in client side
	if job.Mode == model.TaskRunMode_Server {
		return nil
	}
	jobIntersectPath := utils.JobIntersectPath(job.Uuid)
	jobPath := filepath.Dir(jobIntersectPath)
	if err := os.MkdirAll(jobPath, os.ModePerm); err != nil {
		log.Errorf(ctx, "Create job intersect dir error:%v", err)
		return err
	}

	jobFile, err := os.Create(jobIntersectPath)
	if err != nil {
		log.Errorf(ctx, "Create job intersect file error:%v", err)
		return err
	}
	defer jobFile.Close()

	for _, task := range tasks {
		taskResultInfoPath := utils.TaskResultPath(task.Uuid)
		data, err := ioutil.ReadFile(taskResultInfoPath)
		if err != nil {
			log.Errorf(ctx, "Read result file of %s error:%v", task.Uuid, err)
			return err
		}

		data = PostProcessClientResult(data)
		log.Debugf(ctx, "Write from %s to %s (%d) bytes", task.Uuid, job.Uuid, len(data))

		if _, err := jobFile.Write(data); err != nil {
			log.Errorf(ctx, "Write file %s to error:%v", jobIntersectPath, err)
		}
	}

	UpdateJobIntersectCount(ctx, job)
	return nil
}

func UpdateJobIntersectCount(ctx context.Context, job *model.Job) {
	if job.Mode == model.TaskRunMode_Server {
		return
	}

	jobIntersectPath := utils.JobIntersectPath(job.Uuid)

	lineCount, _, _, err := utils.FileMetaInfo(jobIntersectPath)
	if err != nil {
		log.Errorf(ctx, "Get intersect line count of %s error:%v", jobIntersectPath, err)
	}

	job.IntersectCount = lineCount
	if err := model.UpdateJob(job); err != nil {
		log.Errorf(ctx, "Update job intersect file %s line count failed:%v", jobIntersectPath, err)
	}
}

func RemoveDuplicates(ctx context.Context, inputFilePath string) (string, error) {
	log.Infof(ctx, "psi remove duplicates data of file: %s", inputFilePath)

	//check inputFilePath
	if inputFilePath == "" {
		return "", errors.New("Empty inputFilePath.")
	} else {
		_, err := os.Lstat(inputFilePath)
		if os.IsNotExist(err) {
			return "", errors.New("inputFilePath does not exists!")
		}
	}

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return "", errors.New("Failed to open inputFilePath")
	}

	defer inputFile.Close()

	set := make(map[string]struct{}, 10000000)
	var result string

	scanner := bufio.NewScanner(inputFile)
	var lineStr string
	for scanner.Scan() {
		lineStr = scanner.Text()
		if _, ok := set[lineStr]; ok {
			continue
		}

		result += fmt.Sprintf("%s\n", lineStr)

		set[lineStr] = struct{}{}
	}

	outputFilePath := inputFilePath + "_AfterRemoveDuplicates"
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Errorf(ctx, "Create duplicates file:%s,error:%v", outputFilePath, err)
		return outputFilePath, err
	}
	_, err = io.Copy(outputFile, strings.NewReader(result))
	if err != nil {
		log.Errorf(ctx, "Copy in duplicates file:%s,error:%v", outputFilePath, err)
		return outputFilePath, err
	}
	defer outputFile.Close()

	return outputFilePath, nil
}

const (
	PostProcessIntersecthere = false
)

func PostProcessClientResult(data []byte) []byte {
	if data == nil || len(data) == 0 {
		return make([]byte, 0)
	}

	if !PostProcessIntersecthere {
		return data
	}

	s := string(data)
	sArray := strings.SplitN(s, "\n", 3)
	if len(sArray) < 3 {
		return make([]byte, 0)
	}

	return bytes.NewBufferString(sArray[2]).Bytes()
}
