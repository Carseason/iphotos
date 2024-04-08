package iphotos

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type Photos struct {
	values map[string]*Photo
	search Searcher[*SearchItem, []*SearchItem]
	path   string
	mx     *sync.RWMutex
	loger  *slog.Logger
}

// 存储目录
func NewPhotos(p1 string) (*Photos, error) {
	p1, err := filepath.Abs(p1)
	if err != nil {
		return nil, err
	}
	if info, err := os.Stat(p1); err != nil {
		//判断文件夹是否存在
		if os.IsNotExist(err) {
			//创建该文件夹
			err = os.MkdirAll(p1, os.ModePerm)
		}
		if err != nil {
			return nil, err
		}
	} else if !info.IsDir() {
		return nil, errors.New("path is no dir")
	}
	s, err := NewSearch(p1, IndexPropertys, IndexSorts)
	if err != nil {
		return nil, err
	}
	ipos := &Photos{
		path:   p1,
		values: make(map[string]*Photo),
		mx:     &sync.RWMutex{},
		search: s,
		loger:  slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	return ipos, nil
}

// 覆盖日志接口
func (ps *Photos) SetLoger(loger *slog.Logger) error {
	if loger == nil {
		return errors.New("not found loger")
	}
	ps.loger = loger
	return nil
}

// 创建被 photo 调用的context
func (p *Photos) newContext() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		Context:       ctx,
		cancel:        cancel,
		ContextSearch: p.search,
	}
}

// 添加相册
func (ps *Photos) AddPhoto(serialId, path string, excludeDirs []string, excludePaths []string) error {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	// 如果已经存在了
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

// 删除相册
func (ps *Photos) RemovePhoto(serialId string, deleteIndex bool) error {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	p, ok := ps.values[serialId]
	if !ok {
		return errors.New("serialId not found photo")
	}
	// 结束实例
	p.close()
	// 删除实例
	delete(ps.values, serialId)
	// 删除索引
	if deleteIndex {
		// 每次查询100条来删除
		for i := 0; ; i++ {
			datas, err := ps.search.Ids(RequestSearch{
				Limit: 100,
				Filters: map[string]interface{}{
					Index_SerialId: serialId,
				},
			})
			if err != nil || len(datas) == 0 {
				break
			}
			ps.search.Delete(datas...)
		}
	}
	return nil
}

// 查询相册
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

// 只暴露部分搜索接口给业务
func (ps *Photos) Searcher() ContextSearch[*SearchItem, []*SearchItem] {
	return ps.search
}
