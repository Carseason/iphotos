package iphotos

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/sqids/sqids-go"
)

type Photos struct {
	values map[string]*Photo
	search Searcher[*SearchItem, []*SearchItem]
	path   string
	mx     *sync.RWMutex
	sid    *sqids.Sqids
	loger  *slog.Logger
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
	sid, err := sqids.New(sqids.Options{
		Alphabet:  alphabetString,
		MinLength: 16,
	})
	if err != nil {
		return nil, err
	}
	ipos := &Photos{
		path:   p1,
		values: make(map[string]*Photo),
		mx:     &sync.RWMutex{},
		search: s,
		sid:    sid,
		loger:  slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	return ipos, nil
}
func (p *Photos) newContext() *Context {
	return &Context{
		Context:       context.Background(),
		GenFileID:     p.GenFileID,
		ContextSearch: p.search,
	}
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

// 删除相册
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
			datas, err := ps.search.Ids(RequestSearch{
				Limit: 30,
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

// 查询是否索引中
func (ps *Photos) QueryPhotoIndexing(serialId string) bool {
	if p, ok := ps.values[serialId]; ok {
		return p.Indexing()
	}
	return false
}

// serialId,path
func (ps *Photos) GenFileID(vs ...string) (string, error) {
	n := len(vs)
	fids := make([]uint64, 0, n)
	for i := 0; i < n; i++ {
		fids = append(fids, GenFnvID(vs[i]))
	}
	return ps.sid.Encode(fids)
}

// 只暴露部分搜索接口给业务
func (ps *Photos) Searcher() ContextSearch[*SearchItem, []*SearchItem] {
	return ps.search
}

func (ps *Photos) SetLoger(loger *slog.Logger) error {
	if loger == nil {
		return errors.New("not found loger")
	}
	ps.loger = loger
	return nil
}

// 复制数据
func (ps *Photos) CopyData(fids []string, serialId string) error {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	datas, err := ps.search.Query(RequestSearch{
		Ids: fids,
	})
	if err != nil {
		return err
	}
	items := make(map[string]*SearchItem)
	for i := 0; i < len(datas.Result); i++ {
		v := datas.Result[i]
		v.SerialId = serialId
		v.ID = ""
		id, err := ps.GenFileID(serialId, v.Path)
		if err != nil {
			return err
		}
		items[id] = v
	}
	ps.search.Add(items)
	return nil
}
