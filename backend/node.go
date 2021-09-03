package backend

import (
	"smart_proxy/enums"
	//	"github.com/uber-go/atomic"
)

//	Operation on backend node
type BackendNodeOperator struct {
	Target BackendNode
	Op     enums.BACKEND_OPTION
}

// mark backend status
type BackendNodeStatus struct {
	Target     BackendNode
	DownStatus bool
}

type BackendNode struct {
	//State    atomic.Bool //https://pkg.go.dev/go.uber.org/atomic
	sync.RWMutex
	Addr     string
	State    bool   //true - up，false - down
	Metadata string //metadata元数据（JSON String）
}
