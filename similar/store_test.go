package similar

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/rivo/duplo"
)

func TestStore(t *testing.T) {
	s := NewStore("../datas/")

	abspath, err := filepath.Abs("../tmps/20211203_004846.jpg")
	if err != nil {
		t.Error(err)
	}
	hash, _, err := CreateHash(abspath)
	if err != nil {
		t.Error(err)
	}
	s.Adds(map[string]duplo.Hash{
		abspath: *hash,
	})
	if err := s.Save(); err != nil {
		t.Error(err)
	}
}
func TestStoreRegister(t *testing.T) {
	s := NewStore("../datas/")
	if err := s.register(); err != nil {
		t.Error(err)
	}
	abspath, err := filepath.Abs("../tmps/20211203_004847.jpg")
	if err != nil {
		t.Error(err)
	}
	if s.Has(abspath) {
		fmt.Println("已存在")
	}
	datas := s.Query(abspath)
	fmt.Println(datas)
}

func TestStoreQuery(t *testing.T) {
	s := NewStore("")
	abspath, err := filepath.Abs("../tmps")
	if err != nil {
		t.Error(err)
	}
	files, err := os.ReadDir(abspath)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(files); i++ {
		filename := filepath.Join(abspath, files[i].Name())
		hash, _, err := CreateHash(filename)
		if err != nil {
			t.Error(err)
			continue
		}
		s.Add(files[i].Name(), *hash)
	}
	hash, _, err := CreateHash(filepath.Join(abspath, "3243364363955576423_3243364358016232593.jpg"))
	if err != nil {
		t.Error(err)
	}
	datas := s.store.Query(*hash)
	sort.Sort(datas)
	for i := 0; i < len(datas); i++ {
		if datas[i].Score < 60 && datas[i].RatioDiff <= 0.1 {
			fmt.Printf("id: %v, 评分:%v, 绝对差:%v, 汉明距离:%v \n", datas[i].ID, datas[i].Score, datas[i].RatioDiff, datas[i].DHashDistance)
		}
		// 越低匹配度越高
		// fmt.Println(datas[i].Score <= 100 || datas[i].RatioDiff < 0.1)
	}

}

func TestHashsQuery(t *testing.T) {
	s := NewStore("")
	abspath, err := filepath.Abs("../tmps")
	if err != nil {
		t.Error(err)
	}
	files, err := os.ReadDir(abspath)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(files); i++ {
		filename := filepath.Join(abspath, files[i].Name())
		hash, _, err := CreateHash(filename)
		if err != nil {
			t.Error(err)
			continue
		}
		s.Add(files[i].Name(), *hash)
	}

	datas := s.Query(filepath.Join(abspath, "_6080037815858607041_1215f33bb5be7ab6f00.jpg"))
	for i := 0; i < len(datas); i++ {
		fmt.Println(datas[i])
	}

}
