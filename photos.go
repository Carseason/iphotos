package iphotos

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Photos struct {
	values    map[string]*Photo
	Store     *Storer
	Search    Searcher[*SearchItem, []*SearchItem]
	path      string
	mx        *sync.RWMutex
	coverPath string
}
type Context struct {
	context.Context
	Search    Searcher[*SearchItem, []*SearchItem]
	Store     *Storer
	coverPath string
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
	searchPath := filepath.Join(p1, "search")
	storePath := filepath.Join(p1, "duplos")
	coversPath := filepath.Join(p1, "covers")
	if err := mkdirAll(searchPath); err != nil {
		return nil, err
	}
	if err := mkdirAll(storePath); err != nil {
		return nil, err
	}
	if err := mkdirAll(coversPath); err != nil {
		return nil, err
	}
	s, err := NewSearch(searchPath, IndexPropertys, IndexSorts)
	if err != nil {
		return nil, err
	}
	s1, err := NewStore(storePath)
	if err != nil {
		return nil, err
	}
	ipos := &Photos{
		path:      p1,
		values:    make(map[string]*Photo),
		mx:        &sync.RWMutex{},
		Search:    s,
		Store:     s1,
		coverPath: coversPath,
	}
	return ipos, nil
}
func (p *Photos) newContext() *Context {
	return &Context{
		Context:   context.Background(),
		Search:    p.Search,
		Store:     p.Store,
		coverPath: p.coverPath,
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
	if err := ps.Store.Reload(); err != nil {
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
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	if info, err := os.Lstat(path); err != nil { //判断文件夹是否存在
		return err
	} else if !info.IsDir() {
		return errors.New("not found dir")
	}
	if err := p.addPath(path); err != nil {
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
	p.close()
	delete(ps.values, serialId)
	return nil
}
func (ps *Photos) QueryPhoto(serialId string) (*Photo, bool) {
	ps.mx.RLock()
	defer ps.mx.RUnlock()
	p, ok := ps.values[serialId]
	return p, ok
}
func (ps *Photos) Cover(filePath string) string {
	return filepath.Join(ps.coverPath, GenCoverFilename(filePath))
}
