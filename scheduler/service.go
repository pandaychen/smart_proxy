package scheduler

import (
	"github.com/pandaychen/smart_proxy/backend"
	"go.uber.org/zap"
)

var GlobalScheduler *SmartProxyScheduler

// smartproxy controller
type SmartProxyScheduler struct {

	/* Backend operation channel */
	BackendChan chan backend.BackendNodeOperator
	Logger      *zap.Logger
}

func NewSmartProxyScheduler(logger *zap.Logger) (*SmartProxyScheduler, error) {
	sch := &SmartProxyScheduler{
		BackendChan: make(chan backend.BackendNodeOperator),
		Logger:      logger,
	}
	return sch, nil
}

func (s *SmartProxyScheduler) SchedulerLoopRun() {
	s.Logger.Info("Starting SchedulerLoopRun")

	//updates and manages backend nodes
	go func() {
		for {
			select {
			// handle backend operation
			case backend := <-s.BackendChan:
				s.Logger.Info("SchedulerLoopRun handle backendChan", zap.Any("backend", backend))
				//
			}
		}
	}()
}
