package face

import (
	_ "embed"
	"errors"
	"fmt"
	_ "image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"

	pigo "github.com/esimov/pigo/core"
)

//go:embed cascade/facefinder
var cascadeFile []byte

var (
	classifier                    *pigo.Pigo
	Error_ClassifierUninitialized = errors.New("pigo uninitialized")
)

func init() {
	var err error
	p := pigo.NewPigo()
	classifier, err = p.Unpack(cascadeFile)
	if err != nil {
		slog.Error(fmt.Errorf("face init: %s", err).Error())
	}
}

func queryImageFace(inputPath string) ([]pigo.Detection, error) {
	if classifier == nil {
		return nil, Error_ClassifierUninitialized
	}
	src, err := pigo.GetImage(inputPath)
	if err != nil {
		return nil, err
	}
	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y
	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.

	angle := 0.0 // cascade rotation angle. 0.0 is 0 radians and 1.0 is 2*pi radians

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.2)
	return dets, nil
}

// 是否存在人脸
func IsFace(inputPath string) (bool, error) {
	dets, err := queryImageFace(inputPath)
	if err != nil {
		return false, err
	}
	n := len(dets)
	return n > 0, nil
}
