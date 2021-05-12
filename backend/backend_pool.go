package backend

import (
	"github.com/uber-go/atomic"
)

type BackendNode struct {
	State    atomic.Bool //https://pkg.go.dev/go.uber.org/atomic
	Addr     string
	Metadata string //metadata元数据
}

// BackendNodePool is a collection of backend node lists
type BackendNodePool interface {
	//return current pool size
	Size() int

	// get a usable backend nodes from pool
	Get() *BackendNode

	// add a backend node to pool
	Add(addr string)

	// remove a node from pool
	Remove(addr string)

	//
	UpNode(addr string)

	//
	DownNode(addr string)
}
