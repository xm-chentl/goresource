package goresource

import (
	"fmt"

	"github.com/xm-chentl/goresource/dbtype"
)

// Configure 配置 todo: 暂时不封装
type Configure struct {
	Alias   string
	Type    dbtype.Value
	Factory IFactory
}

type factory struct {
	dbTypeOfResource map[dbtype.Value]IResource
	nameOfResource   map[string]IResource
}

func (f factory) BuildByType(dbType dbtype.Value) (resource IResource, err error) {
	resource, ok := f.dbTypeOfResource[dbType]
	if !ok {
		err = fmt.Errorf("dbtype: %s database implementation not configured", dbType.String())
		return
	}

	return
}

func (f factory) BuildByName(name string) (resource IResource, err error) {
	resource, ok := f.nameOfResource[name]
	if !ok {
		err = fmt.Errorf("name: %s database implementation not configured", name)
		return
	}

	return
}

func NewByType(dbTypeConfig map[dbtype.Value]IResource) IFactory {
	return &factory{
		dbTypeOfResource: dbTypeConfig,
		nameOfResource:   make(map[string]IResource),
	}
}

func NewByName(aliasConfig map[string]IResource) IFactory {
	return &factory{
		dbTypeOfResource: make(map[dbtype.Value]IResource),
		nameOfResource:   aliasConfig,
	}
}
