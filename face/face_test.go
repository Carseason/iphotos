package face

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestIsFaces(t *testing.T) {
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

func TestIsFace(t *testing.T) {
	paths := []string{
		filepath.Join("../tmps", "james-lee-0YQz7M2fcYY-unsplash.jpg"),
		filepath.Join("../tmps", "_6080037815858607041_1215f33bb5be7ab6f00.jpg"), // 1
		filepath.Join("../tmps", "_6080037815858607042_121f26d86edfd55b5e0.jpg"), //1
		filepath.Join("../tmps", "1231026887071281154_1.jpg"),
		filepath.Join("../tmps", "martin-martz-qzfu2K5Iz7I-unsplash.jpg"),
		filepath.Join("../tmps", "1591674174665240576_1.jpg"),
		filepath.Join("../tmps", "1653225002987298816_1.jpg"),
	}
	for i := 0; i < len(paths); i++ {
		fmt.Println(IsFace(paths[i]))
	}

}
