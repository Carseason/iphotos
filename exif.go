package iphotos

import (
	"reflect"

	goexif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

// 获取图片exif
func GetImageExif(p1 string) (*PhotoImageExif, error) {
	rawExif, err := goexif.SearchFileAndExtractExif(p1)
	if err != nil {
		return nil, err
	}
	entries, _, err := goexif.GetFlatExifData(rawExif, nil)
	if err != nil {
		return nil, err
	}
	photo := &PhotoImageExif{}
	getGps := func() (*goexif.GpsInfo, error) {
		im, err := exifcommon.NewIfdMappingWithStandard()
		if err != nil {
			return nil, err
		}
		ti := goexif.NewTagIndex()
		_, index, err := goexif.Collect(im, ti, rawExif)
		if err != nil {
			return nil, err
		}
		ifd, err := index.RootIfd.ChildWithIfdPath(exifcommon.IfdGpsInfoStandardIfdIdentity)
		if err != nil {
			return nil, err
		}
		return ifd.GpsInfo()
	}
	elem := reflect.ValueOf(photo).Elem()
	for _, v := range entries {
		elem.FieldByName(v.TagName).Set(reflect.ValueOf(v.Value))
	}
	if gi, err := getGps(); err == nil {
		photo.ExifGps = ExifGps{
			Latitude:  gi.Latitude.Decimal(),
			Longitude: gi.Longitude.Decimal(),
		}
	}
	return photo, nil
}
