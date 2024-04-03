package iphotos

import (
	"fmt"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/rivo/duplo"
)

func TestStore(t *testing.T) {
	s, err := NewStore("./datas/")
	if err != nil {
		t.Error(err)
	}
	p := "tmps/20211203_004846.jpg"
	if img, err := imaging.Open(p); err == nil {
		hash, _ := duplo.CreateHash(img)
		s.Add(p, hash)
	}
	// if err := Store.SaveLocalFile(); err != nil {
	// 	t.Error(err)
	// }
}
func TestStoreRegister(t *testing.T) {
	s, err := NewStore("./datas/")
	if err != nil {
		t.Error(err)
	}
	if err := s.register(); err != nil {
		t.Error(err)
	}
	p := "tmps/20211203_004847.jpg"
	if img, err := imaging.Open(p); err == nil {
		hash, _ := duplo.CreateHash(img)
		ms := s.Query(hash)
		for i := 0; i < ms.Len(); i++ {
			fmt.Println(ms[i].ID)
		}
	}
}
