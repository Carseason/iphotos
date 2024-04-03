package iphotos

type Searcher[T SearchT, TS SearchTS] interface {
	Exist(string) bool
	Query(RequestSearch) (*ResponseSearch[TS], error)
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
	Sorts   []string
	// 查询全部
	All bool
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
	SerialId string `json:"serialId,omitempty"` //序号
	Filename string `json:"filename,omitempty"`
	Path     string `json:"path,omitempty"`
	Size     int64  `json:"size,omitempty"`
	LastDate string `json:"lastDate,omitempty"` //日期 yy-mm-dd
	Tags     any    `json:"tags,omitempty"`     //用户自定义标签
	FileType string `json:"fileType,omitempty"` //文件类型
	//exif
	ExifModel        string `json:"exifModel,omitempty"` //型号
	ExifWidth        int64  `json:"exifWidth,omitempty"`
	ExifHeight       int64  `json:"exifHeight,omitempty"`
	ExifLength       int64  `json:"exifLength,omitempty"`
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
	Index_SerialId = "serialId"
	Index_Filename = "filename"
	Index_FileType = "fileType"
	Index_Tags     = "tags"
	Index_LastDate = "lastDate"
)

var (
	// 用于索引
	IndexPropertys = []string{
		Index_SerialId,
		Index_Filename,
		Index_FileType,
		Index_Tags,
		Index_LastDate,
	}
	// 用于排序
	IndexSorts = []string{
		Index_SerialId,
		Index_Filename,
		Index_LastDate,
	}
)
