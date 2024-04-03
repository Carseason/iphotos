package iphotos

import (
	"encoding/base64"
	"strings"
)

// base64
func toBase64(value string) string {
	return base64.URLEncoding.EncodeToString([]byte(value))
}

// 通过路径生成文件id
func GenFileID(vs ...string) string {
	return strings.Join(vs, "_")
}
