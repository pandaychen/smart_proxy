package pool

import (
	"errors"
	"sync"
)

var ERROR_PARAM_ERROR = errors.New("Pool len illegal")

type TaskInput struct {
	Id        string
	InputData interface{}
}

type TaskResult struct {
	Id         string //if need 回调ID
	OutputData interface{}
	Err        error
	Boolret    bool
}

type SPool struct {
	sync.RWMutex
	WorkNum  int
	Callback func(task TaskInput) (outputData interface{}, err error)
	Retchan  chan TaskResult //结果输出
}

//create
func NewSPool(work_num int, cb func(TaskInput) (interface{}, error)) *SPool {
	if cb == nil {
		return nil
	}
	if work_num <= 0 {
		work_num = 10
	}
	return &SPool{
		Callback: cb,
		WorkNum:  work_num,
		Retchan:  make(chan TaskResult, 1024),
	}
}

func (p *SPool) GetChanResult() chan TaskResult {
	return p.Retchan
}

//todo: 输入优化为chan方式
func (p *SPool) PoolWorkers(TaskInputList []TaskInput) error {
	taskLen := len(TaskInputList)
	if taskLen == 0 {
		return ERROR_PARAM_ERROR
	}
	chInputDataList := make(chan TaskInput, taskLen)
	for _, task := range TaskInputList {
		chInputDataList <- task
	}
	for i := 0; i < p.WorkNum; i++ {
		go func() {
			for {
				p.Lock()
				if len(chInputDataList) == 0 {
					p.Unlock()
					break
				}
				taskData := <-chInputDataList
				p.Unlock()
				//id, outputData, err := p.Callback(taskData)
				outputData, err := p.Callback(taskData)
				p.Retchan <- TaskResult{
					//Id:         id,
					Err:        err,
					OutputData: outputData}
			}
		}()
	}
	return nil
}

//强行关闭任务池
func (p *SPool) Close() {
	close(p.Retchan)
}
