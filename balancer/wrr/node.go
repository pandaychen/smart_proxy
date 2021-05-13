package wrr

import (
	"github.com/pandaychen/smart_proxy/backend"
	"github.com/uber-go/atomic"
)

type BackendNodeWrapper struct {
	backend.BackendNode
	InitWeight      int //初始化权重
	CurrentWeight   int
	EffectiveWeight int //每次pick之后的更新的权重值
}

func CreateWrrBackendNode(addr string, weight int) *BackendNodeWrapper {
	return &BackendNodeWrapper{
		backend.BackendNode{
			State: *atomic.NewBool(true),
			Addr:  addr,
		},
		InitWeight:      weight,
		CurrentWeight:   0,
		EffectiveWeight: weight,
	}
}
