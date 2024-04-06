package iphotos

import "context"

type Context struct {
	context.Context
	ContextSearch[*SearchItem, []*SearchItem]
	GenFileID func(...string) (string, error)
}

type ContextSearch[T SearchT, TS SearchTS] interface {
	Exist(string) bool
	Add(map[string]T) error
	Delete(...string) error
	Query(RequestSearch) (*ResponseSearch[TS], error)
	Ids(RequestSearch) ([]string, error)
}
