package postgres

import "errors"

var (
	ErrUpdateSetArgsIsNotArray   = errors.New("args [0] is not []string for update fields")
	ErrUpdateSetQueryIsNotString = errors.New("args [1] is not string for query condition")
)
