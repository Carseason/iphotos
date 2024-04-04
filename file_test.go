package iphotos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilePath(t *testing.T) {
	fmt.Println(filepath.Dir("./datas"))
	fmt.Println(filepath.Base("./datas"))

	relA := "./datas"
	newA := "./cover"
	relA, err := filepath.Abs(relA)
	if err != nil {
		t.Error(err)
	}
	newA, err = filepath.Abs(newA)
	if err != nil {
		t.Error(err)
	}
	v, err := filepath.Rel(relA, newA)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(v, strings.Contains(newA, relA))
}

func TestFileReadDir(t *testing.T) {
	relA := "../"
	relA, err := filepath.Abs(relA)
	if err != nil {
		t.Error(err)
	}
	files, err := os.ReadDir(relA)
	if err != nil {
		t.Error(err)
	}
	n := len(files)
	for i := 0; i < n; i++ {
		fmt.Println(filepath.Join(relA, files[i].Name()))
	}
}

func TestFilereadDir(t *testing.T) {
	files, err := readDir("./")
	if err != nil {
		t.Error(err)
	}
	n := len(files)
	for i := 0; i < n; i++ {
		fmt.Println(files[i].Name())
	}
}
