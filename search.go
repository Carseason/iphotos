package iphotos

// 必须要先创建该实例才能使用 photo
// 存储目录,用于存储索引的目录
// 包会在该目录下创建索引文件
func NewSearch(p1 string, propertys, sorts []string) (Searcher[*SearchItem, []*SearchItem], error) {
	return NewBleve(
		p1,
		propertys,
		sorts,
	)
}
