package iphotos

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
)

type Bleve[T SearchT, TS SearchTS] struct {
	index     bleve.Index
	absPath   string   //存储目录
	propertys []string //要索引的字段
	sorts     []string //参与排序的字段
	mx        sync.RWMutex
}

func NewBleve[T SearchT, TS SearchTS](path string, propertys, sorts []string) (*Bleve[T, TS], error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	b := &Bleve[T, TS]{
		absPath:   absPath,
		propertys: propertys,
		sorts:     sorts,
		mx:        sync.RWMutex{},
	}
	err = b.createIndex()
	return b, err
}

func (b *Bleve[T, TS]) createIndex() error {
	var err error
	// 索引文件已存在
	if _, err = os.Stat(filepath.Join(b.absPath, "index_meta.json")); err == nil {
		b.index, err = bleve.Open(b.absPath)
	} else {
		// buildMapping 配置功能
		indexMapping := func() mapping.IndexMapping {
			fileStatMapping := bleve.NewDocumentMapping()
			for i := range b.propertys {
				v := bleve.NewTextFieldMapping()
				fileStatMapping.AddFieldMappingsAt(b.propertys[i], v)
			}
			mapping := bleve.NewIndexMapping()
			mapping.DefaultMapping = fileStatMapping
			return mapping
		}()
		b.index, err = bleve.New(b.absPath, indexMapping)
	}
	return err
}

// 删除索引
func (b *Bleve[T, TS]) removeIndex() error {
	if err := b.index.Close(); err != nil {
		return err
	}
	// 如果不存在文件,则认为索引已经不存在
	if _, err := os.Stat(filepath.Join(b.absPath, "index_meta.json")); err != nil {
		return nil
	}
	if err := os.RemoveAll(b.absPath); err != nil {
		return err
	}
	return nil
}

// 重建索引
func (b *Bleve[T, TS]) Reload() error {
	b.mx.Lock()
	defer b.mx.Unlock()
	if err := b.removeIndex(); err != nil {
		return err
	}
	if err := b.createIndex(); err != nil {
		return err
	}
	return nil
}

// 关闭搜索
func (b *Bleve[T, TS]) Close() error {
	b.mx.Lock()
	defer b.mx.Unlock()
	return b.index.Close()
}
func (b *Bleve[T, TS]) Add(values map[string]T) error {
	b.mx.Lock()
	defer b.mx.Unlock()
	for k, v := range values {
		if err := b.index.Index(k, v); err != nil {
			return err
		}
	}
	return nil
}
func (b *Bleve[T, TS]) Delete(keys ...string) error {
	b.mx.Lock()
	defer b.mx.Unlock()
	for _, k := range keys {
		if err := b.index.Delete(k); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bleve[T, TS]) genQuery(req RequestSearch) *bleve.SearchRequest {
	var q []query.Query
	keyword := strings.TrimSpace(req.Keyword)
	if len(keyword) == 0 {
		q = append(q, bleve.NewMatchAllQuery()) //获取全部
	} else {
		req := bleve.NewMatchQuery(keyword)
		q = append(q, req) //根据关键词搜索
	}
	if n := len(req.Filters); n > 0 {
		for k, v := range req.Filters {
			req := bleve.NewMatchQuery(fmt.Sprintf("+%v:'%v'", k, v))
			req.SetField(k)
			q = append(q, req)
		}
	}

	qs := query.NewConjunctionQuery(q)
	return bleve.NewSearchRequestOptions(qs, int(req.Limit), int(req.Offset), req.Explain) //搜索模板，数量，开始，正反排序
}
func (b *Bleve[T, TS]) Query(req RequestSearch) (*ResponseSearch[TS], error) {
	b.mx.RLock()
	defer b.mx.RUnlock()
	searchReq := b.genQuery(req)
	// 排序
	if len(req.Sorts) > 0 {
		var sorts []string
		order := ""
		if !req.Explain {
			order = "-"
		}
		for _, v := range req.Sorts {
			sorts = append(sorts, order+v)
		}
		searchReq.SortBy(sorts)
	}
	searchReq.Fields = []string{"*"} //返回所有结果字段
	res, err := b.index.Search(searchReq)
	if err != nil {
		return nil, err
	}
	//
	count := len(res.Hits)
	if count == 0 {
		return &ResponseSearch[TS]{
			Total: 0,
		}, nil
	}
	datas := make([]map[string]any, count)
	for i := 0; i < count; i++ {
		entrie := res.Hits[i].Fields
		datas[i] = entrie
	}
	var result TS
	by, err := json.Marshal(datas)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(by, &result); err != nil {
		return nil, err
	}
	response := ResponseSearch[TS]{
		Result: result,
		Total:  int64(res.Total),
	}
	return &response, err
}
func (b *Bleve[T, TS]) Exist(id string) bool {
	b.mx.RLock()
	defer b.mx.RUnlock()
	qry := bleve.NewDocIDQuery([]string{id})
	req := bleve.NewSearchRequest(qry)
	// req.Fields = []string{"*"} //返回所有结果字段
	resp, err := b.index.Search(req)
	if err != nil {
		return false
	}
	return len(resp.Hits) > 0
}
func (b *Bleve[T, TS]) Ids(req RequestSearch) ([]string, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()
	searchReq := b.genQuery(req)
	res, err := b.index.Search(searchReq)
	if err != nil {
		return nil, err
	}
	count := len(res.Hits)
	datas := make([]string, 0, count)
	for i := 0; i < count; i++ {
		datas = append(datas, res.Hits[i].ID)
	}
	return datas, nil
}

func (b *Bleve[T, TS]) Count() (int64, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()
	searchReq := b.genQuery(RequestSearch{
		Limit: 0,
	})
	res, err := b.index.Search(searchReq)
	if err != nil {
		return 0, err
	}
	return int64(res.Total), nil
}
