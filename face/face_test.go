package face

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestIsFace(t *testing.T) {
	files, err := os.ReadDir("./tmps")
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(files); i++ {
		fpath := filepath.Join("./tmps", files[i].Name())
		ok, err := IsFace(fpath)
		fmt.Println(files[i].Name(), ":", ok, err)
	}

}
