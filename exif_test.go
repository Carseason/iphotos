package iphotos

import (
	"fmt"
	"testing"
)

func TestGetExif(t *testing.T) {
	e, err := GetImageExif("./tmps/IMG_20240411_160754.jpg")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(
		e.GetImageWidth(),
		e.GetImageHeight(),
		e.GetMake(),
	)
	fmt.Println(e.ExifGps)

}
