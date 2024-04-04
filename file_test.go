package iphotos

import (
	"fmt"
	"io"
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
	relA := "./"
	relA, err := filepath.Abs(relA)
	if err != nil {
		t.Error(err)
	}
	files, err := readDir(relA)
	if err != nil {
		t.Error(err)
	}
	n := len(files)
	for i := 0; i < n; i++ {
		fmt.Println(filepath.Join(relA, files[i].Name()))
	}
}

func TestFileTreeSize(t *testing.T) {
	relA := "./"
	relA, err := filepath.Abs(relA)
	if err != nil {
		t.Error(err)
	}
	f, err := os.Open(relA)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	for {
		files, err := f.Readdir(1)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
		}
		for _, fi := range files {
			fmt.Println(filepath.Join(relA, fi.Name()))
		}
	}
	if err = f.Close(); err != nil {
		t.Error(err)
	}
	if err = f.Close(); err != nil {
		t.Error(err)
	}
}
