package iphotos

import (
	"fmt"
	"testing"
)

func TestPhotos(t *testing.T) {
	p, err := NewPhotos("./datas")
	if err != nil {
		t.Error(err)
		return
	}
	err = p.AddPhoto("1", "./tmps", []string{}, []string{})
	if err != nil {
		t.Error(err)
		return
	}
	datas, err := p.Search.Query(RequestSearch{
		Limit: 10,
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(datas.Result); i++ {
		fmt.Println(datas.Result[i])
	}
	if err := p.Store.Save(); err != nil {
		t.Error(err)
	}
	paths := p.Store.QuerySimilars("./tmps/_6080037815858607041_1215f33bb5be7ab6f00.jpg")
	fmt.Println(paths)
}

func TestPhotosIds(t *testing.T) {
	p, err := NewPhotos("./datas")
	if err != nil {
		t.Error(err)
		return
	}
	datas, err := p.Search.Ids(RequestSearch{
		All: true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(datas); i++ {
		fmt.Println(datas[i])
	}
}
