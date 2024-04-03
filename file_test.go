package iphotos

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestFilePath(t *testing.T) {
	fmt.Println(filepath.Dir("./datas"))
	fmt.Println(filepath.Base("./datas"))
}
