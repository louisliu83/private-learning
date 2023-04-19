package audit

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// FileAuditor autit information into a file
type FileAuditor struct {
	auditFile string
	logger    *logrus.Logger
}

var _ Auditor = &FileAuditor{}

// NewFileAuditor return auditor
func NewFileAuditor(fileName string) *FileAuditor {
	baseDir := filepath.Dir(fileName)
	if err := os.MkdirAll(baseDir, 0666); err != nil {
		logrus.Errorln(err)
		return nil
	}
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	logger := logrus.New()
	logger.SetOutput(file)
	fa := &FileAuditor{
		auditFile: fileName,
		logger:    logger,
	}
	return fa
}

// Audit audit the message into a file
func (a *FileAuditor) Audit(who string, action string, result string, message string, t time.Time) error {
	a.logger.
		WithField("Subject", who).
		WithField("Action", action).
		WithField("Result", result).
		WithTime(t).
		Println(message)
	return nil
}
