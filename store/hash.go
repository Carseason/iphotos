package store

import (
	"image"

	"github.com/disintegration/imaging"
	"github.com/rivo/duplo"
)

func CreateHash(filePath string) (*duplo.Hash, image.Image, error) {
	img, err := imaging.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	hash, img := duplo.CreateHash(img)
	return &hash, img, nil
}
