package goresource

import (
	"fmt"

	"github.com/xm-chentl/goresource/dbtype"
)

type unitOfWork struct {
	uowMap map[dbtype.Value]IUnitOfWork
}

func (u unitOfWork) Commit() (err error) {
	for dbType, uow := range u.uowMap {
		if err = uow.Commit(); err != nil {
			err = fmt.Errorf("[%s] database transaction failed: %s", dbType.String(), err.Error())
			return
		}
	}

	return
}

func Uow() IUnitOfWork {
	return &unitOfWork{
		uowMap: make(map[dbtype.Value]IUnitOfWork),
	}
}
