package iphotos

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/rivo/duplo"
)

// 存储hash
// 用来查找相似照片
type Storer struct {
	*duplo.Store
	_p  string
	_mx *sync.RWMutex
}

func NewStore(p1 string) (*Storer, error) {
	s := &Storer{
		Store: duplo.New(),
		_p:    filepath.Join(p1, "index"),
		_mx:   &sync.RWMutex{},
	}
	s.register()
	return s, nil
}
func (s *Storer) register() error {
	s._mx.Lock()
	defer s._mx.Unlock()
	if _, err := os.Stat(s._p); err != nil {
		return err
	}
	f, err := os.Open(s._p)
	if err != nil {
		return err
	}
	defer f.Close()
	by, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return s.GobDecode(by)
}

// 把数据保存到本地
func (s *Storer) Save() error {
	s._mx.Lock()
	defer s._mx.Unlock()
	if _, err := os.Stat(s._p); err != nil {
		dir, _ := filepath.Split(s._p)
		if err := mkdirAll(dir); err != nil {
			return err
		}
	}
	by, err := s.GobEncode()
	if err != nil {
		return err
	}
	err = os.WriteFile(s._p, by, 0644)
	return err
}
func (s *Storer) Reload() error {
	s._mx.Lock()
	defer s._mx.Unlock()
	// 如果不存在文件,则认为索引已经不存在
	if _, err := os.Stat(s._p); err != nil {
		return nil
	}
	if err := os.RemoveAll(s._p); err != nil {
		return err
	}
	return s.GobDecode([]byte{})
}

// 查询相似照片
func (s *Storer) QuerySimilars(filePath string) []string {
	var paths []string
	if img, err := imaging.Open(filePath); err == nil {
		if hash, img := duplo.CreateHash(img); img != nil {
			datas := s.Query(hash)
			for i := 0; i < datas.Len(); i++ {
				if datas[i].RatioDiff < 0.1 || datas[i].Score <= 0 {
					paths = append(paths, fmt.Sprintf("%v", datas[i].ID))
				}
			}
		}
	}
	return paths
}
