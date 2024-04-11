package iphotos

import (
	"fmt"
	"testing"
	"time"
)

func TestPhotos(t *testing.T) {
	p, err := NewPhotos("./datas")
	if err != nil {
		t.Error(err)
		return
	}
	err = p.AddPhoto("1", "./tmps", nil)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(60 * time.Second)
	datas, err := p.Searcher().Query(RequestSearch{
		Limit: 10,
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(datas.Result); i++ {
		fmt.Println(datas.Result[i])
	}
}

func TestPhotosIds(t *testing.T) {
	p, err := NewPhotos("./datas")
	if err != nil {
		t.Error(err)
		return
	}
	datas, err := p.Searcher().Ids(RequestSearch{})
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
	err = p.AddPhoto("1", "./tmps", nil)
	if err != nil {
		t.Error(err)
		return
	}
	datas, err := p.Searcher().Query(RequestSearch{
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
	err = p.AddPhoto("1", "./tmps", nil)
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
	datas, err := p.Searcher().Query(RequestSearch{
		Limit: 100,
		Filters: map[string]interface{}{
			Index_SerialId: "1",
		},
		Longitude: 33.333,
		Latitude:  44.444,
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(datas.Result); i++ {
		fmt.Println(datas.Result[i].Location)
	}
}
