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
	err = p.AddPhoto("1", "./", []string{}, []string{})
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
	datas, err := p.Search.Ids(RequestSearch{})
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(datas); i++ {
		fmt.Println(datas[i])
	}
}

func TestPhotosSorts(t *testing.T) {
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
		Limit: 100,
		Sorts: []string{
			Index_LastTimestamp,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(datas.Result); i++ {
		fmt.Println(datas.Result[i].LastTimestamp)
	}
}

func TestPhotosRemovePhoto(t *testing.T) {
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
	err = p.RemovePhoto("1", true)
	if err != nil {
		t.Error(err)
		return
	}
}
func TestPhotosQuery(t *testing.T) {
	p, err := NewPhotos("./datas")
	if err != nil {
		t.Error(err)
		return
	}
	datas, err := p.Search.Query(RequestSearch{
		Limit: 10,
		Filters: map[string]interface{}{
			Index_SerialId: "1",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(datas.Result); i++ {
		fmt.Println(datas.Result[i])
	}
}
