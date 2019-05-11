package registry

import (
	"context"
	"fmt"
	"sync"
)

var (
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

type PluginMgr struct {
	plugins map[string]Registry
	lock    sync.Mutex
}

func (p *PluginMgr) registerPlugin(plugin Registry) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	_, ok := p.plugins[plugin.Name()]
	if ok {
		err = fmt.Errorf("duplicate registry plugin")
		return
	}

	p.plugins[plugin.Name()] = plugin
	return
}

func (p *PluginMgr) initRegistry(ctx context.Context, name string,
	opts ...Option) (registry Registry, err error) {
	//查找对应的插件是否存在
	p.lock.Lock()
	defer p.lock.Unlock()
	plugin, ok := p.plugins[name]
	if !ok {
		err = fmt.Errorf("plugin %s not exists", name)
		return
	}

	registry = plugin
	err = plugin.Init(ctx, opts...)
	return
}

// 注册插件
func RegisterPlugin(registry Registry) (err error) {
	return pluginMgr.registerPlugin(registry)
}

// 初始化注册中心
func InitRegistry(ctx context.Context, name string, opts ...Option) (registry Registry, err error) {
	return pluginMgr.initRegistry(ctx, name, opts...)
}
