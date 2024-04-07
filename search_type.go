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
	Public        string `json:"public,omitempty"`        //如果为 "1" 则是可公开的，否则是被手动隐藏了
	//exif
	ExifModel        string `json:"exifModel,omitempty"` //型号
	ExifMake         string `json:"exifMake,omitempty"`  //品牌
	ExifWidth        string `json:"exifWidth,omitempty"`
	ExifHeight       string `json:"exifHeight,omitempty"`
	ExifOriginalDate string `json:"exifOriginalDate,omitempty"`
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
	Index_Public        = "public"
	Index_Model         = "exifModel"
	Index_Make          = "exifMake"
)
const (
	Public_PUBLIC = "1"
	Public_Hidden = "0"
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
		Index_Public,
		Index_Model,
		Index_Make,
	}
	// 用于排序
	IndexSorts = []string{
		Index_SerialId,
		Index_Filename,
		Index_LastTimestamp,
	}
)
