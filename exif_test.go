package iphotos

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetExif(t *testing.T) {
	files, err := os.ReadDir("./tmps")
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(files); i++ {
		fp := filepath.Join("./tmps", files[i].Name())
		e, err := GetImageExif(fp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(
			e.GetGps(),
		)
	}
}
