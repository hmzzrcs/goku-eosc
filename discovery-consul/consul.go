package discovery_consul

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"time"

	"github.com/eolinker/goku-eosc/discovery"
)

type consul struct {
	id         string
	name       string
	address    []string
	params     map[string]string
	labels     map[string]string
	services   discovery.IServices
	context    context.Context
	cancelFunc context.CancelFunc
}

// Start 开始服务发现
func (c *consul) Start() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	c.context = ctx
	c.cancelFunc = cancelFunc

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
	EXIT:
		for {
			select {
			case <-ctx.Done():
				break EXIT
			case <-ticker.C:
				{
					keys := c.services.AppKeys()
					for _, serviceName := range keys {
						nodeSet, err := c.getNodes(serviceName)
						if err != nil {
							log.Error(err)
							continue
						}

						nodes := make([]discovery.INode, 0, len(nodeSet))
						for _, node := range nodeSet {
							nodes = append(nodes, node)
						}
						c.services.Update(serviceName, nodes)
					}
				}

			}

		}

	}()
	return nil
}

// Reset 重置服务发现配置
func (c *consul) Reset(config interface{}, workers map[eosc.RequireId]interface{}) error {
	workerConfig, ok := config.(*Config)
	if !ok {
		return fmt.Errorf("need %s,now %s:%w", eosc.TypeNameOf((*Config)(nil)), eosc.TypeNameOf(config), eosc.ErrorStructType)
	}

	c.address = workerConfig.Config.Address
	c.params = workerConfig.Config.Params
	c.labels = workerConfig.Labels

	return nil
}

// Stop 停止服务发现
func (c *consul) Stop() error {
	c.cancelFunc()
	return nil
}

func (c *consul) Remove(id string) error {
	return c.services.Remove(id)
}

// GetApp 获取服务发现应用
func (c *consul) GetApp(serviceName string) (discovery.IApp, error) {
	nodes, err := c.getNodes(serviceName)
	if err != nil {
		return nil, err
	}

	app, err := c.Create(serviceName, c.labels, nodes)
	if err != nil {
		return nil, err
	}
	c.services.Set(serviceName, app.Id(), app)
	return app, nil
}

// Create 创建服务发现应用
func (c *consul) Create(serviceName string, attrs map[string]string, nodes map[string]discovery.INode) (discovery.IApp, error) {
	return discovery.NewApp(nil, c, attrs, nodes), nil
}

// Id 返回 worker id
func (n *consul) Id() string {
	return n.id
}

func (n *consul) CheckSkill(skill string) bool {
	return discovery.CheckSkill(skill)
}

//getNodes 通过接入地址获取节点信息
func (c *consul) getNodes(service string) (map[string]discovery.INode, error) {
	//TODO Labels怎么处理

	nodeSet := make(map[string]discovery.INode)

	for _, addr := range c.address {
		if !validAddr(addr) {
			log.Errorf("address:%s is invalid", addr)
			continue
		}
		client, err := getConsulClient(addr, c.params)
		if err != nil {
			log.Error(err)
			continue
		}

		clientNodes := getNodesFromClient(client, service)
		for _, node := range clientNodes {
			if _, has := nodeSet[node.Id()]; !has {
				nodeSet[node.Id()] = node
			}
		}
	}

	return nodeSet, nil
}
