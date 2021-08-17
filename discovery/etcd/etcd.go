package etcd

import (
	"context"

	"smart_proxy/backend"
	"smart_proxy/enums"
	etcdtool "smart_proxy/etcd_tools"
	"smart_proxy/scheduler"

	//	"github.com/uber-go/atomic"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type EtcdDiscoveryClient struct {
	Logger      *zap.Logger
	WatchPrefix string
	Scheduler   *scheduler.SmartProxyScheduler //for events report channel
	Client      *etcdtool.EtcdV3Client
}

func NewEtcdDiscoveryClient(logger *zap.Logger, scheduler *scheduler.SmartProxyScheduler, smartproxy_prefix string) (*EtcdDiscoveryClient, error) {
	client := &EtcdDiscoveryClient{
		Logger:      logger,
		WatchPrefix: smartproxy_prefix,
		Scheduler:   scheduler,
	}

	//TODO: fix etcd config
	conf := etcdtool.DefaultConfig()
	etcdcli, err := etcdtool.NewEtcdV3Client(conf)
	if err != nil {
		client.Logger.Error("NewEtcdDiscoveryClient create etcd client error", zap.String("errmsg", err.Error()))
		return nil, err
	}
	client.Client = etcdcli
	return client, nil
}

func (e *EtcdDiscoveryClient) Run() {

	datalist, err := e.Client.GetKeyPrefixValues(e.Client.Context, e.WatchPrefix)
	if err != nil {
		e.Logger.Error("EtcdDiscoveryClient Run error", zap.String("errmsg", err.Error()))
		return
	}

	for key, value := range datalist {
		backendnode := backend.BackendNodeOperator{
			Target: backend.BackendNode{
				//State:    *atomic.NewBool(true),
				State:    true,
				Addr:     value,
				Metadata: key},
			Op: enums.BACKEND_ADD,
		}
		//send to scheduler's channel
		e.Scheduler.BackendChan <- backendnode
	}

	//TODO: fix etcd watcher
	defer e.Client.Close()
	defer e.Client.Cancel()
	for {
		rch := e.Client.Watch(context.Background(), e.WatchPrefix, clientv3.WithPrefix())
		for resp := range rch {
			for _, ev := range resp.Events {
				e.Logger.Info("EtcdDiscoveryClient get loadbalance events", zap.Any("events", ev))
				var backendnode backend.BackendNodeOperator
				if ev.Type == mvccpb.DELETE {
					backendnode = backend.BackendNodeOperator{
						Target: backend.BackendNode{
							State:    true,
							Addr:     string(ev.Kv.Key),
							Metadata: string(ev.Kv.Value)},
						Op: enums.BACKEND_DEL,
					}
				} else {
					backendnode = backend.BackendNodeOperator{
						Target: backend.BackendNode{
							State:    true,
							Addr:     string(ev.Kv.Key),
							Metadata: string(ev.Kv.Value)},
						Op: enums.BACKEND_ADD,
					}
				}
				e.Scheduler.BackendChan <- backendnode
			}
		}
	}
}
