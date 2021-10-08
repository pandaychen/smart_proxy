package wrr

import (
	"errors"
	"smart_proxy/backend"
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

	DownNodeCount int //计数器
}

func NewWrrBalancerPool(logger *zap.Logger, backends_map map[string]int) (*WrrBalancerPool, error) {
	pool := &WrrBalancerPool{
		Logger:          logger,
		BackendNodeSet:  new(sync.Map),
		BackendNodeList: make([]*WrrBackendNodeWrapper, 0),
	}
	for addr, weight := range backends_map {
		pool.AddNode(addr, weight)
	}
	return pool, nil
}

func (p *WrrBalancerPool) Name() string {
	return "weight-rr"
}

func (p *WrrBalancerPool) Size() int {
	//p.RLock()		//fix bugs
	//defer p.RUnlock()
	return len(p.BackendNodeList)
}

// 向 Pool 中添加一个后端节点
func (p *WrrBalancerPool) AddNode(addr string, weight int) error {
	p.Lock()
	defer p.Unlock()
	_, exists := p.BackendNodeSet.LoadOrStore(addr, struct{}{})
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
		return
	}

	p.BackendNodeSet.Delete(addr)

	index := p.index(addr)
	if index >= 0 && index < p.Size() {
		//index 合法
		if !p.BackendNodeList[index].Node.State.IsTrue() {
			//计数器更新
		}
		//remove BackendNodeList[index]
		p.BackendNodeList = append(p.BackendNodeList[:index], p.BackendNodeList[index+1:]...)
	}
}

// 根据 wrr 算法选择后端
// 注意：每次操作单位为 1
//https://pandaychen.github.io/2019/12/15/NGINX-SMOOTH-WEIGHT-ROUNDROBIN-ANALYSIS/#0x04-nginx 平滑的基于权重轮询算法
func (p *WrrBalancerPool) Pick(pick_key string) (*backend.BackendNode, error) {
	var (
		chosen *WrrBackendNodeWrapper
		total  int
	)
	p.RLock()
	defer p.RUnlock()

	for _, node := range p.BackendNodeList {
		//添加状态过滤
		if !node.Node.State.IsTrue() {
			continue
		}
		//lock
		node.Node.Lock()
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
		node.Node.Unlock()
	}

	if chosen != nil {
		chosen.Node.Lock()
		// 更新选中节点的权重值
		chosen.CurrentWeight -= total
		chosen.Node.Unlock()
		return chosen.Node, nil
	}
	return nil, errors.New("None Node Selected")
}

// 查找服务地址在 p.BackendNodeList 中的 index，为了操作对应的 BackendNode
func (p *WrrBalancerPool) index(addr string) int {
	for i, v := range p.BackendNodeList {
		if v.Node.Addr == addr {
			return i
		}
	}
	return -1
}

func (p *WrrBalancerPool) UpNodeStatus(addr string) {
	p.setBackendStatus(addr, true)
	return
}

//关闭后端节点
func (p *WrrBalancerPool) DownNodeStatus(addr string) {
	p.setBackendStatus(addr, false)
	return
}

//设置后端节点状态
func (p *WrrBalancerPool) setBackendStatus(addr string, isNodeUp bool) error {
	//add read lock
	p.RLock()
	node_index := p.index(addr)
	p.RUnlock()

	if node_index < 0 || node_index >= p.Size() {
		return errors.New("illegal index")
	}

	node := p.BackendNodeList[node_index]

	if node.Node.State.IsTrue() {
		if !isNodeUp {
			//down
			node.Node.State.Set(false)
		}
	} else {
		if isNodeUp {
			//up
			node.Node.State.Set(true)
		}
	}

	return nil
}

//返回所有节点，bad节点在先
func (p *WrrBalancerPool) GetAllNodes() []string {
	var (
		bad_iplist  []string
		good_iplist []string
	)
	p.RLock()
	defer p.RUnlock()

	for _, v := range p.BackendNodeList {
		if v.Node.State.IsTrue() {
			good_iplist = append(good_iplist, v.Node.Addr)
		} else {
			bad_iplist = append(bad_iplist, v.Node.Addr)
		}
	}

	bad_iplist = append(bad_iplist, good_iplist...)
	return bad_iplist
}
