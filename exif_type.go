package iphotos

import (
	"strconv"

	"github.com/qichengzx/coordtransform"
)

type PhotoImageExif struct {

	// 光圈值
	ApertureValue any `json:"ApertureValue,omitempty"`

	// 亮度值
	BrightnessValue any `json:"BrightnessValue,omitempty"`

	// 色彩空间
	ColorSpace any `json:"ColorSpace,omitempty"`

	// 组件配置
	ComponentsConfiguration any `json:"ComponentsConfiguration,omitempty"`

	// 压缩
	Compression any `json:"Compression,omitempty"`

	// 日期时间
	DateTime any `json:"DateTime,omitempty"`

	// 日期时间数字化
	DateTimeDigitized any `json:"DateTimeDigitized,omitempty"`

	// 日期时间原始
	DateTimeOriginal any `json:"DateTimeOriginal,omitempty"`

	// EXIF标签
	ExifTag any `json:"ExifTag,omitempty"`

	// Exif版本
	ExifVersion any `json:"ExifVersion,omitempty"`

	// 曝光偏差值
	ExposureBiasValue any `json:"ExposureBiasValue,omitempty"`

	// 曝光模式
	ExposureMode any `json:"ExposureMode,omitempty"`

	// 曝光计划
	ExposureProgram any `json:"ExposureProgram,omitempty"`

	// 曝光时间
	ExposureTime any `json:"ExposureTime,omitempty"`

	// F编号
	FNumber any `json:"FNumber,omitempty"`

	// 闪光
	Flash any `json:"Flash,omitempty"`

	// Flashpix版本
	FlashpixVersion any `json:"FlashpixVersion,omitempty"`

	// 焦距
	FocalLength any `json:"FocalLength,omitempty"`

	// 焦距35mm胶片
	FocalLengthIn35mmFilm any `json:"FocalLengthIn35mmFilm,omitempty"`

	// GPS海拔
	GPSAltitude any `json:"GPSAltitude,omitempty"`

	// GPS海拔参考
	GPSAltitudeRef any `json:"GPSAltitudeRef,omitempty"`

	// GPS日期戳
	GPSDateStamp any `json:"GPSDateStamp,omitempty"`

	// GPS纬度
	GPSLatitude any `json:"GPSLatitude,omitempty"`

	// PS纬度参考
	GPSLatitudeRef any `json:"GPSLatitudeRef,omitempty"`

	// GPS经度
	GPSLongitude any `json:"GPSLongitude,omitempty"`

	// GPS经度参考
	GPSLongitudeRef any `json:"GPSLongitudeRef,omitempty"`

	// GPSProcessingMethod
	GPSProcessingMethod any `json:"GPSProcessingMethod,omitempty"`

	// GPS标签
	GPSTag any `json:"GPSTag,omitempty"`

	// GPS时间戳
	GPSTimeStamp any `json:"GPSTimeStamp,omitempty"`

	// iso速度等级
	ISOSpeedRatings any `json:"ISOSpeedRatings,omitempty"`

	// 图片高度
	ImageLength any `json:"ImageLength,omitempty"`

	// 图片宽度
	ImageWidth any `json:"ImageWidth,omitempty"`

	// 互操作性指数
	InteroperabilityIndex any `json:"InteroperabilityIndex,omitempty"`

	// 互操作性标签
	InteroperabilityTag any `json:"InteroperabilityTag,omitempty"`

	// 互操作性版本
	InteroperabilityVersion any `json:"InteroperabilityVersion,omitempty"`

	// JPEG交换格式
	JPEGInterchangeFormat any `json:"JPEGInterchangeFormat,omitempty"`

	// JPEGInterchangeFormatLength
	JPEGInterchangeFormatLength any `json:"JPEGInterchangeFormatLength,omitempty"`

	// 光源
	LightSource any `json:"LightSource,omitempty"`

	// 品牌
	Make any `json:"Make,omitempty"`

	// 最大光圈值
	MaxApertureValue any `json:"MaxApertureValue,omitempty"`

	// 测光模式
	MeteringMode any `json:"MeteringMode,omitempty"`

	// 手机型号
	Model any `json:"Model,omitempty"`

	// 方向
	Orientation any `json:"Orientation,omitempty"`

	// 像素X尺寸
	PixelXDimension any `json:"PixelXDimension,omitempty"`

	// 像素Y尺寸
	PixelYDimension any `json:"PixelYDimension,omitempty"`

	// 分辨率单位
	ResolutionUnit any `json:"ResolutionUnit,omitempty"`

	// 场景捕获类型
	SceneCaptureType any `json:"SceneCaptureType,omitempty"`

	// 传感方法
	SensingMethod any `json:"SensingMethod,omitempty"`

	// 快门速度值
	ShutterSpeedValue any `json:"ShutterSpeedValue,omitempty"`

	// 亚秒时间
	SubSecTime any `json:"SubSecTime,omitempty"`

	// 亚秒时间数字化
	SubSecTimeDigitized any `json:"SubSecTimeDigitized,omitempty"`

	// 亚秒时间原始
	SubSecTimeOriginal any `json:"SubSecTimeOriginal,omitempty"`

	// 白平衡
	WhiteBalance any `json:"WhiteBalance,omitempty"`

	// X分辨率
	XResolution any `json:"XResolution,omitempty"`

	// YCBCR定位
	YCbCrPositioning any `json:"YCbCrPositioning,omitempty"`

	// Y分辨率
	YResolution any `json:"YResolution,omitempty"`
	ExifGps
}
type ExifGps struct {
	// 经度
	Longitude float64 `json:"Longitude,omitempty"`
	// 纬度
	Latitude float64 `json:"Latitude,omitempty"`
}

// 获取经纬度
func (p *ExifGps) GetGps() (float64, float64) {
	return p.Longitude, p.Latitude
}

// WGS84toGCJ02 WGS84坐标系->火星坐标系
func (p *ExifGps) GetGpsWGS84() (float64, float64) {
	return coordtransform.WGS84toGCJ02(p.GetGps())
}
func (p *ExifGps) GetGpsString() string {
	return strconv.FormatFloat(p.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(p.Latitude, 'f', -1, 64)
}

// 获取宽度
func (p *PhotoImageExif) GetImageWidth() int64 {
	return anyToInt64Value(p.ImageWidth)
}

// 获取高度
func (p *PhotoImageExif) GetImageHeight() int64 {
	return anyToInt64Value(p.ImageLength)
}

// 获取手机型号
func (p *PhotoImageExif) GetMake() string {
	if v, ok := p.Make.(string); ok {
		return v
	}
	return ""
}

// 获取品牌
func (p *PhotoImageExif) GetModel() string {
	if v, ok := p.Model.(string); ok {
		return v
	}
	return ""
}

// 获取原始日期
func (p *PhotoImageExif) GetDateTimeOriginal() string {
	if v, ok := p.DateTimeOriginal.(string); ok {
		return v
	}
	return ""
}

// 获取iso等级
func (p *PhotoImageExif) GetISOSpeedRatings() int64 {
	return anyToInt64Value(p.ISOSpeedRatings)
}

// 获取曝光计划
func (p *PhotoImageExif) GetExposureProgram() int64 {
	return anyToInt64Value(p.ExposureProgram)
}
