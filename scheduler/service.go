package scheduler

import (
	"context"
	"fmt"
	"smart_proxy/backend"
	"smart_proxy/enums"
	"smart_proxy/reverseproxy"
	"time"

	"go.uber.org/zap"
)

var GlobalScheduler *SmartProxyScheduler

// smartproxy controller
type SmartProxyScheduler struct {
	/* Backend operation channel */
	BackendChan  chan backend.BackendNodeOperator
	PeerStatChan chan backend.PeerStateOperator
	Logger       *zap.Logger
	ReverseGroup *reverseproxy.SmartReverseProxyGroup

	//for statistics without lock
	PeerStatus map[string]map[string]int
	PeerMap    map[string]string //temp
}

func NewSmartProxyScheduler(logger *zap.Logger, reverse_group *reverseproxy.SmartReverseProxyGroup, dis2schChan chan backend.BackendNodeOperator, peerStatChan chan backend.PeerStateOperator) (*SmartProxyScheduler, error) {
	sch := &SmartProxyScheduler{
		Logger:       logger,
		ReverseGroup: reverse_group,
		BackendChan:  dis2schChan,
		PeerStatChan: peerStatChan,
		PeerStatus:   make(map[string]map[string]int),
		PeerMap:      make(map[string]string),
	}
	return sch, nil
}

func (s *SmartProxyScheduler) Run(ctx context.Context) {
	s.Logger.Info("Starting SchedulerLoopRun")

	statTicker := time.NewTicker(10 * time.Second)

	//updates and manages backend nodes
	go func() {
		for {
			select {
			// handle backend operation
			case backend := <-s.BackendChan:
				s.Logger.Info("SchedulerLoopRun handle backendChan", zap.Any("backend", backend))
				s.ProcessBackendNodes(&backend)
			case peerstat := <-s.PeerStatChan:
				s.Logger.Info("SchedulerLoopRun handle peerStatChan", zap.Any("backend", peerstat))
				switch peerstat.Op {
				case enums.BACKEND_BAD:
					if _, exists := s.PeerStatus[peerstat.Target.Addr]; !exists {
						s.PeerStatus[peerstat.Target.Addr] = make(map[string]int)
					}
					s.PeerMap[peerstat.Target.Addr] = peerstat.Target.ProxyName
					s.PeerStatus[peerstat.Target.Addr][string(enums.BACKEND_BAD)]++
					s.PeerStatus[peerstat.Target.Addr][string(enums.BACKEND_TOTAL)]++
				case enums.BACKEND_GOOD:
					if _, exists := s.PeerStatus[peerstat.Target.Addr]; !exists {
						s.PeerStatus[peerstat.Target.Addr] = make(map[string]int)
					}
					s.PeerMap[peerstat.Target.Addr] = peerstat.Target.ProxyName
					s.PeerStatus[peerstat.Target.Addr][string(enums.BACKEND_TOTAL)]++
				default:
					continue
				}
			case <-statTicker.C:
				s.Logger.Info("SchedulerLoopRun handle statTicker")
				for peer, tmap := range s.PeerStatus {
					s.Logger.Info("SchedulerLoopRun handle statTicker", zap.String("peer", peer), zap.Any("map", tmap))
					var bad, total int
					var proxyname string
					if _, exists := tmap[string(enums.BACKEND_BAD)]; !exists {
						continue
					} else {
						bad = tmap[string(enums.BACKEND_BAD)]
					}
					if _, exists := tmap[string(enums.BACKEND_TOTAL)]; !exists {
						continue
					} else {
						total = tmap[string(enums.BACKEND_TOTAL)]
					}
					proxyname = s.PeerMap[peer]
					fmt.Println(bad / total)
					err_rate := float64(1 / float64(2))
					if err_rate > 0.3 {
						//Down peers
						s.UpDownBackendNodes(proxyname, peer, enums.BACKEND_DOWN)
					}
					//reset map
					s.PeerStatus = make(map[string]map[string]int)
					s.PeerMap = make(map[string]string)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

//设置backend节点状态
func (s *SmartProxyScheduler) UpDownBackendNodes(proxy_name string, backend_addr string, op enums.BACKEND_OPTION) error {
	switch op {
	case enums.BACKEND_DOWN:
		proxy := s.ReverseGroup.GetReverseProxy(proxy_name)
		if proxy != nil {
			proxy.BackendNodePool.DownNodeStatus(backend_addr)
		}
	case enums.BACKEND_UP:
		proxy := s.ReverseGroup.GetReverseProxy(proxy_name)
		if proxy != nil {
			proxy.BackendNodePool.UpNodeStatus(backend_addr)
		}
	default:
		return nil
	}
	return nil
}

func (s *SmartProxyScheduler) ProcessBackendNodes(node *backend.BackendNodeOperator) {
	switch node.Op {
	case enums.BACKEND_DEL:
		proxy := s.ReverseGroup.GetReverseProxy(node.Target.ProxyName)
		if proxy != nil {
			proxy.BackendNodePool.RemoveNode(node.Target.Addr)
		}
	case enums.BACKEND_ADD:
		proxy := s.ReverseGroup.GetReverseProxy(node.Target.ProxyName)
		if proxy != nil {
			proxy.BackendNodePool.AddNode(node.Target.Addr, 1)
		}
	default:
		return
	}
}
