package backend

import (
	"github.com/uber-go/atomic"
)

//	Operation on backend
type BackendOption struct {
	target BackendNode
	op     enums.BACKEND_OPTION
}

type BackendNode struct {
	State    atomic.Bool //https://pkg.go.dev/go.uber.org/atomic
	Addr     string
	Metadata string //metadata元数据
}

// BackendNodePool is a collection of backend node lists（with properly load balance choice）
type BackendNodePool interface {
	//return current pool size
	Size() int

	// get a usable backend nodes from pool
	Get() *BackendNode

	// add a backend node to pool
	Add(addr string)

	// remove a node from pool
	Remove(addr string)

	// set a node up status
	UpNode(addr string)

	// set a node down status
	DownNode(addr string)
}
