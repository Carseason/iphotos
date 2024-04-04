package iphotos

import (
	"strconv"
	"strings"

	goexif "github.com/dsoprea/go-exif/v3"
)

type Exif struct {
	LastDate         string `json:"exifLastDate,omitempty"`
	ExifModel        string `json:"exifModel,omitempty"` //型号
	ExifWidth        string `json:"exifWidth,omitempty"`
	ExifHeight       string `json:"exifHeight,omitempty"`
	ExifLength       string `json:"exifLength,omitempty"`
	ExifOriginalDate string `json:"exifOriginalDate,omitempty"`
}

func GetImageExifData(p1 string) ([]goexif.ExifTag, *goexif.MiscellaneousExifData, error) {
	rawExif, err := goexif.SearchFileAndExtractExif(p1)
	if err != nil {
		return nil, nil, err
	}
	return goexif.GetFlatExifData(rawExif, &goexif.ScanOptions{})

}
func GetImageExif(p1 string) (*Exif, error) {
	tags, _, err := GetImageExifData(p1)
	if err != nil {
		return nil, err
	}
	var item Exif
	for i := 0; i < len(tags); i++ {
		switch strings.TrimSpace(tags[i].TagName) {
		case "ImageWidth":
			if v, ok := anyToInt64(tags[i].Value); ok {
				item.ExifWidth = strconv.FormatInt(v, 10)
			}
		case "ImageHeight":
			if v, ok := anyToInt64(tags[i].Value); ok {
				item.ExifHeight = strconv.FormatInt(v, 10)
			}
		case "Model":
			if v, ok := tags[i].Value.(string); ok {
				item.ExifModel = v
			}
		case "ImageLength":
			if v, ok := anyToInt64(tags[i].Value); ok {
				item.ExifLength = strconv.FormatInt(v, 10)
			}
		case "DateTime":
			if v, ok := tags[i].Value.(string); ok {
				item.LastDate = v
			}
		case "DateTimeOriginal":
			if v, ok := tags[i].Value.(string); ok {
				item.ExifOriginalDate = v
			}
		}
		// fmt.Printf("%v: %v\n", tags[i].TagName, tags[i].Value)
	}
	return &item, nil
}
