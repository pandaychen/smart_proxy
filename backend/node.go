package backend

import (
	"smart_proxy/enums"
	"sync"

	//	"github.com/uber-go/atomic"
	atom "smart_proxy/pkg/pyatomic"
)

//	Operation on backend node
type BackendNodeOperator struct {
	Target BackendNode
	Op     enums.BACKEND_OPTION
}

//	Operation on backend node
type PeerStateOperator struct {
	Target BackendNode
	Op     enums.BACKEND_STAT
}

// mark backend status
type BackendNodeStatus struct {
	Target     BackendNode
	DownStatus bool
}

type BackendNode struct {
	//State    atomic.Bool //https://pkg.go.dev/go.uber.org/atomic
	sync.RWMutex
	ProxyName string //belongs to which proxy
	Addr      string
	State     *atom.AtomicBool //true - up，false - down，标识节点状态
	Metainfo  Metadata         //metadata元数据（JSON String）
}

type Metadata struct {
	Addr   string `json:"addr"`
	Weight int    `json:"weight"`
}
