package discovery

import "github.com/eolinker/eosc/internal"

type Services struct {
	apps        internal.IUntyped
	appNameOfId internal.IUntyped
}

func NewServices() *Services {
	return &Services{apps: internal.NewUntyped(), appNameOfId: internal.NewUntyped()}
}

func (n *Services) get(namespace string) (internal.IUntyped, bool) {
	v, ok := n.apps.Get(namespace)
	if !ok {
		return nil, ok
	}
	apps, ok := v.(internal.IUntyped)
	return apps, ok
}

func (s *Services) Set(serviceName string, id string, app IApp) error {
	s.appNameOfId.Set(id, serviceName)
	if apps, ok := s.get(serviceName); ok {
		apps.Set(id, app)
		return nil
	}
	apps := internal.NewUntyped()
	apps.Set(id, app)
	s.apps.Set(serviceName, apps)
	return nil
}

func (s *Services) Remove(id string) error {
	v, has := s.appNameOfId.Del(id)
	if has {
		apps, ok := s.get(v.(string))
		if ok {
			apps.Del(id)
		}
	}
	return nil
}

func (s *Services) Update(serviceName string, nodes []INode) error {
	if apps, ok := s.get(serviceName); ok {
		for _, r := range apps.List() {
			v, ok := r.(IApp)
			if ok {
				v.Reset(nodes)
			}
		}
	}
	return nil
}

func (s *Services) AppKeys() []string {
	return s.apps.Keys()
}

type IServices interface {
	Set(serviceName string, id string, app IApp) error
	Remove(id string) error
	Update(serviceName string, nodes []INode) error
	AppKeys() []string
}
