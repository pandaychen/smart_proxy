package wrr

import (
	"smart_proxy/backend"
	//	"github.com/uber-go/atomic"
)

// WrrBackendNodeWrapper 代表一个后端
type WrrBackendNodeWrapper struct {
	backend.BackendNode     // 通用的后端节点结构
	InitWeight          int // 初始化权重（固定不变）
	CurrentWeight       int // 当前权重（初始值为 0）
	EffectiveWeight     int // 每次 Pick 之后的更新的权重值（初始值等于 InitWeight ）
}

func NewWrrBackendNode(addr string, weight int) *WrrBackendNodeWrapper {
	return &BackendNodeWrapper{
		backend.BackendNode{
			State: true, // 默认初始化为开启状态
			//State: *atomic.NewBool(true),
			Addr: addr,
		},
		InitWeight:      weight,
		CurrentWeight:   0, // 初始化为 0
		EffectiveWeight: weight,
	}
}
