package iphotos

import (
	"context"
	"errors"
)

type Context struct {
	context.Context
	cancel context.CancelFunc
	ContextSearch[*SearchItem, []*SearchItem]
}

type ContextSearch[T SearchT, TS SearchTS] interface {
	Exist(string) bool
	Add(map[string]T) error
	Delete(...string) error
	Query(RequestSearch) (*ResponseSearch[TS], error)
	Ids(RequestSearch) ([]string, error)
}

var (
	ErrContextClose = errors.New("context close")
)
