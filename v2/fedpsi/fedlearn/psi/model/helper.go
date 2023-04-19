package model

// TaskCanStop test if a task in 'taskStatus' can be stopped
func TaskCanStop(taskStatus string) bool {
	return taskStatus == TaskStatus_Running ||
		taskStatus == TaskStatus_ClientWaiting
}

// JobCanStop test if a job in 'jobStatus' can be stopped
func JobCanStop(jobStatus string) bool {
	return jobStatus == JobStatus_Confirmed ||
		jobStatus == JobStatus_Ready ||
		jobStatus == JobStatus_Running ||
		jobStatus == JobStatus_WaitingPartyConfirm
}
