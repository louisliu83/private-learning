package scheduler

import (
	"pa.cn/fedlearn/psi/server/manager"
)

type SharderScheduler struct {
}

var _ Scheduler = NewSharderScheduler()

func NewSharderScheduler() *SharderScheduler {
	return &SharderScheduler{}
}

func (s *SharderScheduler) Tick() {
	// No need impl
}

func (s *SharderScheduler) Start() {
	ctx := contexForScheduler()
	manager.GetSharder().Run(ctx)
}
