package iphotos

import (
	"fmt"
	"testing"
)

func TestGenCover(t *testing.T) {
	p, err := GenCover("./datas/covers", "./tmps/MVIMG_20240316_193529.jpg", 220, 220)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("path:", p)
}
func TestGenVideoCover(t *testing.T) {
	p, err := GenVideoCover("./tmps/test.mp4", "./datas/covers")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("path:", p)
}
