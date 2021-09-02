package scheduler

import (
	"go.uber.org/zap"
	"smart_proxy/backend"
)

var GlobalScheduler *SmartProxyScheduler

// smartproxy controller
type SmartProxyScheduler struct {
	/* Backend operation channel */
	BackendChan chan backend.BackendNodeOperator
	Logger      *zap.Logger
}

func NewSmartProxyScheduler(logger *zap.Logger, dis2schChan chan backend.BackendNodeOperator) (*SmartProxyScheduler, error) {
	sch := &SmartProxyScheduler{
		Logger:      logger,
		BackendChan: dis2schChan,
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
			}
		}
	}()
}
