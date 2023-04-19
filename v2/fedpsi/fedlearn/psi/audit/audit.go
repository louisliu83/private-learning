package audit

import (
	"sync"
	"time"

	"fedlearn/psi/common/config"

	"github.com/sirupsen/logrus"
)

// Auditor will log the audit information
type Auditor interface {
	Audit(who string, action string, result string, message string, t time.Time) error
}

const (
	defaultAuditFile = "/tmp/psi/audit.log"
)

var (
	fileAuditor       Auditor
	fileAuditorLocker = sync.Mutex{}
)

// GetAuditor return an auditor
func GetAuditor() Auditor {
	fileAuditorLocker.Lock()
	defer fileAuditorLocker.Unlock()
	auditorFile := config.GetConfig().Audit.File
	if auditorFile == "" {
		auditorFile = defaultAuditFile
	}
	if fileAuditor == nil {
		fileAuditor = NewFileAuditor(auditorFile)
	}
	return fileAuditor
}

// Log the audit information
func Log(who string, action string, result string, message string) {
	if !config.GetConfig().FeatureGate.Audit {
		return
	}

	a := GetAuditor()
	if a == nil {
		logrus.Errorln("Cannot get the auditor")
	}

	a.Audit(who, action, result, message, time.Now())
}
