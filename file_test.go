package iphotos

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilePath(t *testing.T) {
	fmt.Println(filepath.Dir("./datas"))
	fmt.Println(filepath.Base("./datas"))

	relA := "/Users/carseason./github/go/iphotos/datas"
	newA := "/Users/carseason./github/go/iphotos/datas/cover"
	v, err := filepath.Rel(relA, newA)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(v, strings.Contains(newA, relA))
}
