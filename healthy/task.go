package healthy

import (
	"errors"
	sphttp "smart_proxy/pkg/http"
	grpool "smart_proxy/pkg/pool"
)

type Task struct {
	Name     string
	Addr     string
	CheckRet bool
}

func TcpCheckTask(rawTask grpool.TaskInput) (outputData interface{}, err error) {
	task, ok := rawTask.InputData.(Task)
	if !ok {
		return nil, errors.New("wrong format")
	}

	checkret, err := sphttp.CheckTcpAlive(task.Addr)
	if checkret {
		task.CheckRet = true
		return task, nil
	}
	task.CheckRet = false
	return task, err
}
