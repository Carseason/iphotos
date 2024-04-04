package iphotos

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	s, err := NewSearch("./datas", IndexPropertys, IndexSorts)
	if err != nil {
		t.Error(err)
	}
	err = s.Reload()
	fmt.Println(err)
}

func TestSearchExist(t *testing.T) {
	s, err := NewSearch("./datas", IndexPropertys, IndexSorts)
	if err != nil {
		t.Error(err)
	}
	ok := s.Exist("8:L1VzZXJzL2NhcnNlYXNvbi4vZ2l0aHViL2dvL3Bob3Rvcy90bXBzLzIwMjExMjAzXzAwNDg0Ni5qcGc=")
	fmt.Println(ok)
}

func TestSearchQuery(t *testing.T) {
	s, err := NewSearch("./datas", IndexPropertys, IndexSorts)
	if err != nil {
		t.Error(err)
	}
	val, err := s.Query(RequestSearch{
		// Filters: map[string]interface{}{
		// 	Index_SerialId: "8",
		// },
		// Sorts: []string{
		// 	Index_Filename,
		// },
		Limit: 10,
	})
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(val.Result); i++ {
		fmt.Println(val.Result[i])
	}
}

func TestSearchCount(t *testing.T) {
	s, err := NewSearch("./datas", IndexPropertys, IndexSorts)
	if err != nil {
		t.Error(err)
	}
	total, err := s.Count()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(total)
}
