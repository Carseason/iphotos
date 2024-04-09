package face

import (
	"embed"
	"errors"
	"fmt"
	_ "image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"log/slog"

	pigo "github.com/esimov/pigo/core"
)

//go:embed cascade/facefinder
var cascadeFile []byte

//go:embed cascade/puploc
var puplocFile []byte

//go:embed cascade/lps
var flpcsDir embed.FS

var (
	classifier                    *pigo.Pigo
	plc                           *pigo.PuplocCascade
	Error_ClassifierUninitialized = errors.New("pigo uninitialized")
	flpcs                         map[string][]*pigo.FlpCascade
)
var (
	eyeCascades  = []string{"lp46", "lp44", "lp42", "lp38", "lp312"}
	mouthCascade = []string{"lp93", "lp84", "lp82", "lp81"}
)

func init() {
	var err error
	p := pigo.NewPigo()
	classifier, err = p.Unpack(cascadeFile)
	if err != nil {
		slog.Error(fmt.Errorf("p.Unpack: %s", err).Error())
		return
	}
	plc, err = pigo.NewPuplocCascade().UnpackCascade(puplocFile)
	if err != nil {
		slog.Error(fmt.Errorf("plc.UnpackCascade: %s", err).Error())
		return
	}

	cascades, err := flpcsDir.ReadDir("cascade/lps")
	if err != nil {
		slog.Error(fmt.Errorf("flpcsDir.Open: %s", err).Error())
		return
	}
	lps, err := fs.Sub(flpcsDir, "cascade/lps")
	if err != nil {
		slog.Error(fmt.Errorf("fs.Sub: %s", err).Error())
		return
	}
	flpcs = make(map[string][]*pigo.FlpCascade)
	for _, cascade := range cascades {
		f, err := lps.Open(cascade.Name())
		if err != nil {
			slog.Error(fmt.Errorf("lps.Open: %s", err).Error())
			continue
		}
		defer f.Close()
		by, err := io.ReadAll(f)
		if err != nil {
			slog.Error(fmt.Errorf("io.ReadAll: %s", err).Error())
			continue
		}
		flpc, err := plc.UnpackCascade(by)
		if err != nil {
			slog.Error(fmt.Errorf(" plc.UnpackCascade: %s", err).Error())
			continue
		}
		flpcs[cascade.Name()] = append(flpcs[cascade.Name()], &pigo.FlpCascade{
			PuplocCascade: flpc,
		})
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
		MinSize:     60,
		MaxSize:     600,
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

	// Run the classifier over the obtained leaf nodes and return the detection dets.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.2)
	return dets, nil
}

// 人脸特征检测
func queryImageFaceLandmark(inputPath string) ([][]int, error) {
	if classifier == nil {
		return nil, Error_ClassifierUninitialized
	}
	src, err := pigo.GetImage(inputPath)
	if err != nil {
		return nil, err
	}
	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y
	imgParams := pigo.ImageParams{
		Pixels: pixels,
		Rows:   rows,
		Cols:   cols,
		Dim:    cols,
	}
	cParams := pigo.CascadeParams{
		MinSize:     60,
		MaxSize:     600,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: imgParams,
	}
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.

	angle := 0.0 // cascade rotation angle. 0.0 is 0 radians and 1.0 is 2*pi radians

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	results := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	results = classifier.ClusterDetections(results, 0.2)
	n := len(results)
	dets := make([][]int, n)

	for i := 0; i < n; i++ {
		dets[i] = append(dets[i], results[i].Row, results[i].Col, results[i].Scale, int(results[i].Q), 0)
		// left eye
		puploc := &pigo.Puploc{
			Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
			Col:      results[i].Col - int(0.185*float32(results[i].Scale)),
			Scale:    float32(results[i].Scale) * 0.4,
			Perturbs: 63,
		}
		leftEye := plc.RunDetector(pigo.Puploc{
			Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
			Col:      results[i].Col - int(0.185*float32(results[i].Scale)),
			Scale:    float32(results[i].Scale) * 0.4,
			Perturbs: 63,
		}, imgParams, angle, false)
		// right eye
		puploc = &pigo.Puploc{
			Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
			Col:      results[i].Col + int(0.185*float32(results[i].Scale)),
			Scale:    float32(results[i].Scale) * 0.4,
			Perturbs: 63,
		}
		rightEye := plc.RunDetector(*puploc, imgParams, 0.0, false)
		if rightEye.Row > 0 && rightEye.Col > 0 {
			dets[i] = append(dets[i], rightEye.Row, rightEye.Col, int(rightEye.Scale), int(results[i].Q), 1)
		}
		for _, eye := range eyeCascades {
			for _, flpc := range flpcs[eye] {
				flp := flpc.GetLandmarkPoint(leftEye, rightEye, imgParams, puploc.Perturbs, false)
				if flp.Row > 0 && flp.Col > 0 {
					dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
				}
				flp = flpc.GetLandmarkPoint(leftEye, rightEye, imgParams, puploc.Perturbs, true)
				if flp.Row > 0 && flp.Col > 0 {
					dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
				}
			}
		}
		// Traverse all the mouth cascades and run the detector on each of them.
		for _, mouth := range mouthCascade {
			for _, flpc := range flpcs[mouth] {
				flp := flpc.GetLandmarkPoint(leftEye, rightEye, imgParams, puploc.Perturbs, false)
				if flp.Row > 0 && flp.Col > 0 {
					dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
				}
			}
		}
		flp := flpcs["lp84"][0].GetLandmarkPoint(leftEye, rightEye, imgParams, puploc.Perturbs, true)
		if flp.Row > 0 && flp.Col > 0 {
			dets[i] = append(dets[i], flp.Row, flp.Col, int(flp.Scale), int(results[i].Q), 2)
		}
	}

	return dets, nil
}

// 是否存在人脸
func IsFace(inputPath string) (bool, error) {
	dets, err := queryImageFaceLandmark(inputPath)
	if err != nil {
		return false, err
	}
	n := len(dets)
	if n >= 2 {
		return true, nil
	}
	return false, nil
}
