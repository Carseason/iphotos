package iphotos

import "testing"

func TestGenCover(t *testing.T) {
	if err := GenCover("./datas/covers", "./tmps/MVIMG_20240316_193529.jpg", 220, 220); err != nil {
		t.Error(err)
	}

}
