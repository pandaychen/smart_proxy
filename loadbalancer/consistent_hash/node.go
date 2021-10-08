package consistent_hash

import (
	"smart_proxy/backend"
	atom "smart_proxy/pkg/pyatomic"
)

// ConsistenthashBackendNodeWrapper 代表一个后端
type ConsistenthashBackendNodeWrapper struct {
	Node *backend.BackendNode // 通用的后端节点结构
}

func NewConsistenthashBackendNodeWrapper(addr string) *ConsistenthashBackendNodeWrapper {
	return &ConsistenthashBackendNodeWrapper{
		Node: &backend.BackendNode{
			State: atom.NewAtomicBoolWithVal(true), // 默认初始化为开启状态
			Addr:  addr,
		},
	}
}
