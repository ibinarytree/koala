package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ibinarytree/koala/registry"
)

const (
	MaxServiceNum          = 8
	MaxSyncServiceInterval = time.Second * 10
)

//etcd 注册插件
type EtcdRegistry struct {
	options   *registry.Options
	client    *clientv3.Client
	serviceCh chan *registry.Service

	value              atomic.Value
	lock               sync.Mutex
	registryServiceMap map[string]*RegisterService
}

type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}

type RegisterService struct {
	id          clientv3.LeaseID
	service     *registry.Service
	registered  bool
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

var (
	etcdRegistry *EtcdRegistry = &EtcdRegistry{
		serviceCh:          make(chan *registry.Service, MaxServiceNum),
		registryServiceMap: make(map[string]*RegisterService, MaxServiceNum),
	}
)

func init() {

	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	etcdRegistry.value.Store(allServiceInfo)
	registry.RegisterPlugin(etcdRegistry)
	go etcdRegistry.run()
}

//插件的名字
func (e *EtcdRegistry) Name() string {
	return "etcd"
}

//初始化
func (e *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {

	e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}

	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.options.Addrs,
		DialTimeout: e.options.Timeout,
	})

	if err != nil {
		err = fmt.Errorf("init etcd failed, err:%v", err)
		return
	}

	return
}

//服务注册
func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {

	select {
	case e.serviceCh <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}
	return
}

//服务反注册
func (e *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	return
}

func (e *EtcdRegistry) run() {

	ticker := time.NewTicker(MaxSyncServiceInterval)
	for {
		select {
		case service := <-e.serviceCh:
			registryService, ok := e.registryServiceMap[service.Name]
			if ok {
				for _, node := range service.Nodes {
					registryService.service.Nodes = append(registryService.service.Nodes, node)
				}
				registryService.registered = false
				break
			}
			registryService = &RegisterService{
				service: service,
			}
			e.registryServiceMap[service.Name] = registryService
		case <-ticker.C:
			e.syncServiceFromEtcd()
		default:
			e.registerOrKeepAlive()
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (e *EtcdRegistry) registerOrKeepAlive() {

	for _, registryService := range e.registryServiceMap {
		if registryService.registered {
			e.keepAlive(registryService)
			continue
		}

		e.registerService(registryService)
	}
}

func (e *EtcdRegistry) keepAlive(registryService *RegisterService) {

	select {
	case resp := <-registryService.keepAliveCh:
		if resp == nil {
			registryService.registered = false
			return
		}
	}
	return
}

func (e *EtcdRegistry) registerService(registryService *RegisterService) (err error) {

	resp, err := e.client.Grant(context.TODO(), e.options.HeartBeat)
	if err != nil {
		return
	}

	registryService.id = resp.ID
	for _, node := range registryService.service.Nodes {

		tmp := &registry.Service{
			Name: registryService.service.Name,
			Nodes: []*registry.Node{
				node,
			},
		}

		data, err := json.Marshal(tmp)
		if err != nil {
			continue
		}

		key := e.serviceNodePath(tmp)
		fmt.Printf("register key:%s\n", key)
		_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
		if err != nil {
			continue
		}

		// the key 'foo' will be kept forever
		ch, err := e.client.KeepAlive(context.TODO(), resp.ID)
		if err != nil {
			continue
		}

		registryService.keepAliveCh = ch
		registryService.registered = true
	}

	return
}

func (e *EtcdRegistry) serviceNodePath(service *registry.Service) string {

	nodeIP := fmt.Sprintf("%s:%d", service.Nodes[0].IP, service.Nodes[0].Port)
	return path.Join(e.options.RegistryPath, service.Name, nodeIP)
}

func (e *EtcdRegistry) servicePath(name string) string {
	return path.Join(e.options.RegistryPath, name)
}

func (e *EtcdRegistry) getServiceFromCache(ctx context.Context,
	name string) (service *registry.Service, ok bool) {

	allServiceInfo := e.value.Load().(*AllServiceInfo)
	//一般情况下，都会从缓存中读取
	service, ok = allServiceInfo.serviceMap[name]
	return
}

func (e *EtcdRegistry) GetService(ctx context.Context,
	name string) (service *registry.Service, err error) {

	//一般情况下，都会从缓存中读取
	service, ok := e.getServiceFromCache(ctx, name)
	if ok {
		return
	}

	//如果缓存中没有这个service，则从etcd中读取
	e.lock.Lock()
	defer e.lock.Unlock()
	//先检测，是否已经从etcd中加载成功了
	service, ok = e.getServiceFromCache(ctx, name)
	if ok {
		return
	}

	//从etcd中读取指定服务名字的服务信息
	key := e.servicePath(name)
	resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}

	service = &registry.Service{
		Name: name,
	}

	for _, kv := range resp.Kvs {
		value := kv.Value
		var tmpService registry.Service
		err = json.Unmarshal(value, &tmpService)
		if err != nil {
			return
		}

		for _, node := range tmpService.Nodes {
			service.Nodes = append(service.Nodes, node)
		}
	}

	allServiceInfoOld := e.value.Load().(*AllServiceInfo)
	var allServiceInfoNew = &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	for key, val := range allServiceInfoOld.serviceMap {
		allServiceInfoNew.serviceMap[key] = val
	}

	allServiceInfoNew.serviceMap[name] = service
	e.value.Store(allServiceInfoNew)
	return
}

func (e *EtcdRegistry) syncServiceFromEtcd() {

	var allServiceInfoNew = &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	ctx := context.TODO()
	allServiceInfo := e.value.Load().(*AllServiceInfo)

	//对于缓存的每一个服务，都需要从etcd中进行更新
	for _, service := range allServiceInfo.serviceMap {
		key := e.servicePath(service.Name)
		resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
		if err != nil {
			allServiceInfoNew.serviceMap[service.Name] = service
			continue
		}

		serviceNew := &registry.Service{
			Name: service.Name,
		}

		for _, kv := range resp.Kvs {
			value := kv.Value
			var tmpService registry.Service
			err = json.Unmarshal(value, &tmpService)
			if err != nil {
				fmt.Printf("unmarshal failed, err:%v value:%s", err, string(value))
				return
			}

			for _, node := range tmpService.Nodes {
				serviceNew.Nodes = append(serviceNew.Nodes, node)
			}
		}
		allServiceInfoNew.serviceMap[serviceNew.Name] = serviceNew
	}

	e.value.Store(allServiceInfoNew)
}
