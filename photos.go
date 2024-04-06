package iphotos

import (
	"context"
	"errors"
	"hash/fnv"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sqids/sqids-go"
)

type Photos struct {
	values map[string]*Photo
	search Searcher[*SearchItem, []*SearchItem]
	path   string
	mx     *sync.RWMutex
	sid    *sqids.Sqids
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
		Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
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

func (ps *Photos) GenFileID(vs ...string) (string, error) {
	f := fnv.New64a()
	f.Write([]byte(strings.Join(vs, ":")))
	return ps.sid.Encode([]uint64{
		f.Sum64()},
	)
}

// 为文件添加标签
func (ps *Photos) AddTag(fid string, tag string) error {
	ps.mx.Lock()
	defer ps.mx.Unlock()
	datas, err := ps.search.Query(RequestSearch{
		Filters: map[string]interface{}{
			Index_SerialId: fid,
		},
	})
	if err != nil {
		return err
	}
	items := make(map[string]*SearchItem)
	for i := 0; i < len(datas.Result); i++ {
		isExist := false
		tags := datas.Result[i].GetTags()
		for j := 0; j < len(tags); j++ {
			if tags[i] == tag {
				isExist = true
				break
			}
		}
		if !isExist {
			tags = append(tags, tag)
			datas.Result[i].Tags = tags
			items[fid] = datas.Result[i]
		}
	}
	if len(items) > 0 {
		return ps.search.Add(items)
	}
	return nil
}

// 只暴露部分搜索接口给业务
func (ps *Photos) Searcher() ContextSearch[*SearchItem, []*SearchItem] {
	return ps.search
}
