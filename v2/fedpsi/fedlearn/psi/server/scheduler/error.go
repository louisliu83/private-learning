package scheduler

import (
	"fmt"
)

func NewTaskScheduleError(taskUID string, err error) error {
	return fmt.Errorf("Schedule task %s error:%v", taskUID, err)
}

func NewJobScheduleError(jobUID string, err error) error {
	return fmt.Errorf("Schedule job %s error:%v", jobUID, err)
}
