package iphotos

import (
	"context"
	"errors"

	"github.com/carseason/iphotos/similar"
)

type Context struct {
	context.Context
	cancel   context.CancelFunc
	Searcher ContextSearch[*SearchItem, []*SearchItem]
	Similar  *similar.Storer
}

type ContextSearch[T SearchT, TS SearchTS] interface {
	Exist(string) bool
	Add(map[string]T) error
	Delete(...string) error
	Query(RequestSearch) (*ResponseSearch[TS], error)
	Ids(RequestSearch) ([]string, error)
	Hidden(ids ...string) error
}

var (
	ErrContextClose = errors.New("context close")
)
