package wrr

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

//wrr: https://github.com/pandaychen/goes-wrapper/blob/master/balancer/roundrobin.go
//https://github.com/nginx/nginx/commit/52327e0627f49dbda1e8db695e63a4b0af4448b1

// WrrBalancerPool is a instance of backend.BackendNodePool
type WrrBalancerPool struct {
	sync.RWMutex
	BackendNodeList []*WrrBackendNodeWrapper
	Logger          *zap.Logger
	BackendNodeSet  *sync.Map // 用于去重
}

func NewWrrBalancerPool(logger *zap.Logger, backends_map map[string]int) (*WrrBalancerPool, error) {
	pool := &WrrBalancerPool{
		Logger:          logger,
		BackendsSet:     new(sync.Map),
		BackendNodeList: make([]*WrrBackendNodeWrapper, 0),
	}
	for addr, weight := range backends_map {
		pool.Add(addr, weight)
	}
	return pool, nil
}

func (p *WrrBalancerPool) Name() string {
	return "wrr"
}

func (p *WrrBalancerPool) Size() int {
	p.Rlock()
	defer p.UnRlock()
	return len(p.BackendNodeList)
}

// 向 Pool 中添加一个后端节点
func (p *WrrBalancerPool) AddNode(addr string, weight int) error {
	p.Lock()
	defer p.Unlock()
	_, exists := p.BackendNodeSet.LoadOrStore(addr, struct{})
	if exists {
		return errors.New("Node Exists")
	}

	if weight <= 0 {
		weight = 1
	}

	//Create New Node
	bnode := NewWrrBackendNode(addr, weight)
	p.BackendNodeList = append(p.BackendNodeList, bnode)
	return nil
}

// 向 Pool 中移除指定的后端节点
func (p *WrrBalancerPool) RemoveNode(addr string) {
	p.Lock()
	defer p.Unlock()

	_, exists := p.BackendNodeSet.Load(addr)
	if !exists {
		return errors.New("Node Not Exists")
	}

	p.BackendsSet.Delete(addr)

	index := p.index(addr)
	if index >= 0 && index < p.Size() {
		//index 合法
		if p.BackendNodeList[idx].State == false {

		}
		//remove BackendNodeList[index]
		p.BackendNodeList = append(p.BackendNodeList[:index], p.BackendNodeList[index+1:]...)
	}
}

// 根据 wrr 算法选择后端
// 注意：每次操作单位为 1
//https://pandaychen.github.io/2019/12/15/NGINX-SMOOTH-WEIGHT-ROUNDROBIN-ANALYSIS/#0x04-nginx 平滑的基于权重轮询算法
func (p *WrrBalancerPool) Pick(pick_key string) (*backend.BackendNode, error) {
	p.RLock()
	defer p.RUnlock()

	var (
		chosen *backend.BackendNode
		total  int
	)

	for _, node := range p.BackendNodeList {
		//lock
		node.Lock()
		total += node.EffectiveWeight
		node.CurrentWeight += node.EffectiveWeight

		if node.EffectiveWeight < node.InitWeight {
			node.EffectiveWeight++
		}

		if chosen == nil || node.CurrentWeight > chosen.CurrentWeight {
			//update choice
			chosen = node
		}
		//unlock
		node.Unlock()
	}

	if chosen != nil {
		chosen.Lock()
		// 更新选中节点的权重值
		chosen.CurrentWeight -= total
		chosen.Unlock()
		return chosen, nil
	}
	return nil, errors.New("None Node Selected")
}

// 查找服务地址在 p.BackendNodeList 中的 index，为了操作对应的 BackendNode
func (p *WrrBalancerPool) index(addr string) int {
	for i, v := range p.BackendNodeList {
		if v.Addr == addr {
			return i
		}
	}
	return -1
}

func (p *WrrBalancerPool) UpNodeStatus(addr string) {
	return
}

func (p *WrrBalancerPool) DownNodeStatus(addr string) {
	return
}
