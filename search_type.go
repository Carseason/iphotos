package iphotos

type Searcher[T SearchT, TS SearchTS] interface {
	Exist(string) bool
	Query(RequestSearch) (*ResponseSearch[TS], error)
	Count() (int64, error)
	Ids(RequestSearch) ([]string, error)
	Add(map[string]T) error
	Delete(...string) error
	Close() error
	Reload() error
	Hidden(ids ...string) error
}

// 搜索结果
type ResponseSearch[T SearchTS] struct {
	Result T     `json:"hits"`
	Total  int64 `json:"estimatedTotalHits,omitempty"`
}

// 过滤提交
type Filter struct {
	Key   string
	Value []string
}

// 查询搜索
type RequestSearch struct {
	Keyword string
	Offset  int64
	Limit   int64
	// Filter: "id > 1 AND genres = Action",
	Filters map[string]interface{}
	Ids     []string
	Sorts   []string
	// 倒叙
	Explain bool

	// 经度
	Longitude float64 `json:"Longitude,omitempty"`
	// 纬度
	Latitude float64 `json:"Latitude,omitempty"`
}

type SearchTS interface {
	[]*SearchItem
}
type SearchT interface {
	*SearchItem
}

const (
	FileType_IMAGE = "image"
	FileType_Video = "video"
)

// 入库模型
// bleve 不支持 []int64
// 所以使用 []string 代替
type SearchItem struct {
	ID            string `json:"id,omitempty"`       //搜索id
	SerialId      string `json:"serialId,omitempty"` //序号
	Filename      string `json:"filename,omitempty"`
	Path          string `json:"path,omitempty"`
	Size          string `json:"size,omitempty"`
	LastDate      string `json:"lastDate,omitempty"`      //日期 yy-mm-dd
	LastTimestamp string `json:"lastTimestamp,omitempty"` //时间戳,秒
	Tags          any    `json:"tags,omitempty"`          //用户自定义标签
	FileType      string `json:"fileType,omitempty"`      //文件类型
	FileExt       string `json:"fileExt,omitempty"`       //文件后缀
	Status        string `json:"status,omitempty"`        //1为公开,别的则为隐藏或删除了
	Identify      string `json:"identify,omitempty"`      //图片内容类型,比如 face 为人脸
	//exif
	ExifModel        string `json:"exifModel,omitempty"` //型号
	ExifMake         string `json:"exifMake,omitempty"`  //品牌
	ExifOriginalDate string `json:"exifOriginalDate,omitempty"`
	// 经纬度
	// [Longitude,Latitude]
	Location []float64 `json:"location,omitempty"` // Built up using Latitude and Longitude
}

func (s *SearchItem) IsImage() bool {
	return s.FileType == FileType_IMAGE
}
func (s *SearchItem) IsVideo() bool {
	return s.FileType == FileType_Video
}
func (s *SearchItem) GetTags() []string {
	return anyToStrings(s.Tags)
}

const (
	Index_SerialId      = "serialId"
	Index_Filename      = "filename"
	Index_FileType      = "fileType"
	Index_FileExt       = "fileExt"
	Index_Tags          = "tags"
	Index_LastTimestamp = "lastTimestamp"
	Index_Model         = "exifModel"
	Index_Make          = "exifMake"
	Index_Status        = "status"
	Index_Identify      = "identify"
	Index_Path          = "path"
)
const (
	// 数据状态
	Status_Public = "1" //公开
	Status_Hidden = "0" //移除
	// 数据内容识别
	Identify_Face      = "face"      //人脸
	Identify_Thing     = "thing"     //物品
	Identify_Landscape = "landscape" //风景
)

var (
	// 用于索引
	IndexPropertys = []string{
		Index_SerialId,
		Index_Filename,
		Index_FileType,
		Index_FileExt,
		Index_Tags,
		Index_LastTimestamp,
		Index_Model,
		Index_Make,
		Index_Status,
		Index_Identify,
		Index_Path,
	}
	// 用于排序
	IndexSorts = []string{
		Index_SerialId,
		Index_Filename,
		Index_LastTimestamp,
		Index_Path,
	}
)
