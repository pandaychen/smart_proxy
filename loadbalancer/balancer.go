package loadbalancer

import (
	"smart_proxy/backend"
)

//LoadBalance算法必须实现的接口，实现需要确保协程安全
//	BackendNodePool is a collection of backend node lists（with properly load balance choice）
type BackendNodePool interface {
	//return lb name
	Name() string

	//return current pool size
	Size() int

	// get a usable backend nodes from pool with lb's pick method
	Pick(pick_key string) (*backend.BackendNode, error)

	// add a backend node to pool
	AddNode(addr string, weight int) error

	// remove a node from pool
	RemoveNode(addr string)

	// set a node up status
	UpNodeStatus(addr string)

	// set a node down status
	DownNodeStatus(addr string)
}
