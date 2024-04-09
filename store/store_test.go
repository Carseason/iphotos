package store

import (
	"fmt"
	"path/filepath"
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
