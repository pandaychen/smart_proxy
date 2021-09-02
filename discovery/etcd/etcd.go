package etcd

import (
	"context"

	"smart_proxy/backend"
	"smart_proxy/enums"
	etcdtool "smart_proxy/pkg/etcd_tools"

	//"smart_proxy/scheduler"

	//	"github.com/uber-go/atomic"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// EtcdConfig
type EtcdConfig struct {
	Cluster    string
	RootPrefix string //Root 前缀

	PasswordAuthOn bool
	Username       string
	Password       string
	Cert           string
	Key            string
	CommonName     string
	TrustedCaCert  string
	ResultChan     chan backend.BackendNodeOperator
	Logger         *zap.Logger
}

type EtcdDiscoveryClient struct {
	Logger      *zap.Logger
	WatchPrefix string
	//Scheduler   *scheduler.SmartProxyScheduler //for events report channel
	BackendChan chan backend.BackendNodeOperator
	Client      *etcdtool.EtcdV3Client
}

func NewEtcdDiscoveryClient(etcd_conf *EtcdConfig) (*EtcdDiscoveryClient, error) {
	client := &EtcdDiscoveryClient{
		Logger:      etcd_conf.Logger,
		WatchPrefix: etcd_conf.RootPrefix,
		BackendChan: etcd_conf.ResultChan,
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
		e.BackendChan <- backendnode
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
				e.BackendChan <- backendnode
			}
		}
	}
}

func (e *EtcdDiscoveryClient) Close() {
	//TODO: 通知channel关闭
	e.Client.Close()
	e.Client.Cancel()
}
