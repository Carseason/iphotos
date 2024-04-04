package iphotos

import (
	"fmt"
	"testing"
)

func TestExif(t *testing.T) {
	relA := "./tmps/MVIMG_20240316_193529.jpg"
	row, err := GetImageExif(relA)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(row.ExifModel, row.ExifMake)
}

func TestExifData(t *testing.T) {
	relA := "./tmps/MVIMG_20240316_193529.jpg"
	tags, _, err := GetImageExifData(relA)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(tags); i++ {
		fmt.Printf("%v: %v\n", tags[i].TagName, tags[i].Value)
	}
}
