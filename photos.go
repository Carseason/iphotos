package iphotos

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
)

type Photos struct {
	values map[string]*Photo
	Search Searcher[*SearchItem, []*SearchItem]
	path   string
	mx     *sync.RWMutex
}
type Context struct {
	context.Context
	Search Searcher[*SearchItem, []*SearchItem]
}

// 存储目录
func NewPhotos(p1 string) (*Photos, error) {
	p1, err := filepath.Abs(p1)
	if err != nil {
		return nil, err
	}
	if err := mkdirAll(p1); err != nil {
		return nil, err
	}
	s, err := NewSearch(p1, IndexPropertys, IndexSorts)
	if err != nil {
		return nil, err
	}
	ipos := &Photos{
		path:   p1,
		values: make(map[string]*Photo),
		mx:     &sync.RWMutex{},
		Search: s,
	}
	return ipos, nil
}
func (p *Photos) newContext() *Context {
	return &Context{
		Context: context.Background(),
		Search:  p.Search,
	}
}

// 重建索引
func (ps *Photos) ReloadIndex() error {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	var errs error
	if err := ps.Search.Reload(); err != nil {
		errs = errors.Join(errs, err)
	}
	return errs
}

// 添加相册
func (ps *Photos) AddPhoto(serialId, path string, excludeDirs []string, excludePaths []string) error {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	if _, ok := ps.values[serialId]; ok {
		return errors.New("serialId existed")
	}
	// 排除本身存储的路径
	excludePaths = append(excludePaths, ps.path)
	//
	p, err := NewPhoto(ps.newContext(), serialId, excludeDirs, excludePaths)
	if err != nil {
		return err
	}
	if err := p.AddPATH(path); err != nil {
		return err
	}
	ps.values[serialId] = p
	return nil
}
func (ps *Photos) RemovePhoto(serialId string, deleteIndex bool) error {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	p, ok := ps.values[serialId]
	if !ok {
		return errors.New("serialId not found photo")
	}
	p.close()
	delete(ps.values, serialId)
	// 删除索引
	if deleteIndex {
		// 每次查询30条来删除
		for i := 0; ; i++ {
			datas, err := ps.Search.Ids(RequestSearch{
				Limit: 30,
				Filters: map[string]interface{}{
					Index_SerialId: serialId,
				},
			})
			if err != nil || len(datas) == 0 {
				break
			}
			ps.Search.Delete(datas...)
		}
	}
	return nil
}
func (ps *Photos) QueryPhoto(serialId string) (*Photo, bool) {
	ps.mx.RLock()
	defer ps.mx.RUnlock()
	p, ok := ps.values[serialId]
	return p, ok
}
func (ps *Photos) QueryPhotoIndexing(serialId string) bool {
	if p, ok := ps.values[serialId]; ok {
		return p.Indexing()
	}
	return false
}
