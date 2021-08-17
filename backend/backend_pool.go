package backend

import (
	"smart_proxy/enums"
	//	"github.com/uber-go/atomic"
)

//	Operation on backend
type BackendNodeOperator struct {
	Target BackendNode
	Op     enums.BACKEND_OPTION
}

type BackendNode struct {
	//State    atomic.Bool //https://pkg.go.dev/go.uber.org/atomic
	State    bool
	Addr     string
	Metadata string //metadata元数据
}

// BackendNodePool is a collection of backend node lists（with properly load balance choice）
type BackendNodePool interface {
	//return current pool size
	Size() int

	// get a usable backend nodes from pool
	Pick() *BackendNode

	// add a backend node to pool
	Add(addr string)

	// remove a node from pool
	Remove(addr string)

	// set a node up status
	UpNode(addr string)

	// set a node down status
	DownNode(addr string)
}
