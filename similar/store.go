package similar

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rivo/duplo"
)

// 存储hash
// 用来查找相似照片
type Storer struct {
	ctx    context.Context
	cancel context.CancelFunc
	store  *duplo.Store
	p      string
	mx     *sync.RWMutex
}

// 持久化的路径, 是否注册旧数据
// 是否定时写入数据
func NewStore(path string) *Storer {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Storer{
		store:  duplo.New(),
		mx:     &sync.RWMutex{},
		ctx:    ctx,
		cancel: cancel,
	}
	if len(path) > 0 {
		s.p = filepath.Join(path, "similars")
		s.register()
		go s.sleep()
	}
	return s
}
func (s *Storer) Close() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.cancel()
	return s.store.GobDecode([]byte{})
}
func (s *Storer) register() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if _, err := os.Stat(s.p); err != nil {
		return err
	}
	f, err := os.Open(s.p)
	if err != nil {
		return err
	}
	defer f.Close()
	by, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return s.store.GobDecode(by)
}
func (s *Storer) sleep() {
	for {
		ticker := time.NewTicker(1 * time.Hour)
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if err := s.Save(); err != nil {
				slog.Error(err.Error())
			}
		}
	}
}

// 把数据保存到本地
func (s *Storer) Save() error {
	if err := s.ctx.Err(); err != nil {
		return err
	}
	s.mx.Lock()
	defer s.mx.Unlock()
	if _, err := os.Stat(s.p); err != nil {
		//判断文件夹是否存在
		if os.IsNotExist(err) {
			// 切割路径
			dir, _ := filepath.Split(s.p)
			// 创建文件夹
			err = os.MkdirAll(dir, os.ModePerm) //创建该文件夹
		}
		if err != nil {
			return err
		}
	}
	by, err := s.store.GobEncode()
	if err != nil {
		return err
	}
	return os.WriteFile(s.p, by, 0644)
}

// 清空数据
func (s *Storer) Clear() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if err := s.store.GobDecode([]byte{}); err != nil {
		return err
	}
	// 如果文件存在,则删除该文件
	if _, err := os.Stat(s.p); err == nil {
		return os.Remove(s.p)
	}
	return nil
}

// 查询相似照片
func (s *Storer) Query(filePath string) []string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	var paths []string
	if img, err := imaging.Open(filePath); err == nil {
		if hash, img := duplo.CreateHash(img); img != nil {
			datas := s.store.Query(hash)
			sort.Sort(datas)
			for i := 0; i < len(datas); i++ {
				// 如果 score 为负数，则认为照片几乎一样
				// 如果 score < 60, 则认为照片角度不一定一样, 再计算绝对差来区别
				//
				if datas[i].Score <= 0 {
					if v, ok := datas[i].ID.(string); ok {
						paths = append(paths, v)
					}
				} else if datas[i].Score < 60 && datas[i].RatioDiff <= 0.1 && datas[i].DHashDistance < 10 {
					if v, ok := datas[i].ID.(string); ok {
						paths = append(paths, v)
					}
				}
			}
		}
	}
	return paths
}
func (s *Storer) Adds(values map[string]duplo.Hash) {
	s.mx.Lock()
	defer s.mx.Unlock()
	for k, v := range values {
		s.store.Add(k, v)
	}
}
func (s *Storer) Add(id string, hash duplo.Hash) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.store.Add(id, hash)
}
func (s *Storer) Has(id string) bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.store.Has(id)
}
func (s *Storer) Delete(id string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.store.Delete(id)
}
func (s *Storer) Ids(id string) []string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	datas := s.store.IDs()
	var ids []string
	for i := 0; i < len(datas); i++ {
		if v, ok := datas[i].(string); ok {
			ids = append(ids, v)
		}
	}
	return ids
}
