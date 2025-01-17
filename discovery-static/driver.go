package discovery_static

import (
	"reflect"
	"sync"

	"github.com/eolinker/goku-eosc/discovery"

	"github.com/eolinker/eosc"
)

const (
	driverName = "static"
)

//driver 实现github.com/eolinker/eosc.eosc.IProfessionDriver接口
type driver struct {
	profession string
	name       string
	driver     string
	label      string
	desc       string
	configType reflect.Type
	params     map[string]string
}

func NewDriver() *driver {
	return &driver{configType: reflect.TypeOf(new(Config))}
}

func (d *driver) ConfigType() reflect.Type {
	return d.configType
}

func (d *driver) Create(id, name string, v interface{}, workers map[eosc.RequireId]interface{}) (eosc.IWorker, error) {
	s := &static{
		id:     id,
		name:   name,
		locker: sync.RWMutex{},
		apps:   make(map[string]discovery.IApp),
	}
	s.Reset(v, workers)
	return s, nil
}
