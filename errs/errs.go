package errs

import "errors"

var (
	ResIsNotPtr            = errors.New("res is not ptr")
	ResIsNotSlice          = errors.New("res is not slice")
	ResIsNotStruct         = errors.New("res is not struct")
	ResIsNotIDbModel       = errors.New("res is not an implementation of dbfactory.IDbModel")
	UpdateFullNotAllowed   = errors.New("full update not allowed")
	DeleteFullNotAllowed   = errors.New("full delete not allowed")
	QueryArgsError         = errors.New("args parameter error")
	QueryGrammarEmptyError = errors.New("query grammar empty")
	GrammarError           = errors.New("sql grammar error")
)
