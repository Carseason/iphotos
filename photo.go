package iphotos

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	goexif "github.com/dsoprea/go-exif/v3"
	"github.com/fsnotify/fsnotify"
	"github.com/rivo/duplo"
)

type Photo struct {
	watcher  *fsnotify.Watcher
	serialId string
	// 排除的文件夹名称,在里面的文件夹和子文件夹则排除
	excludeDirs map[string]struct{}
	// 排除的路径
	excludePaths map[string]struct{}
	ctx          *Context
}

// 创建相册程序,如果有多个相册,则需要创建多个
func NewPhoto(ctx *Context, serialId string, excludeDirs, excludePaths []string) (*Photo, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	p := &Photo{
		watcher:      w,
		serialId:     serialId,
		excludeDirs:  map[string]struct{}{},
		excludePaths: map[string]struct{}{},
		ctx:          ctx,
	}
	for _, k := range excludeDirs {
		p.excludeDirs[k] = struct{}{}
	}
	for _, k := range excludePaths {
		p.excludePaths[k] = struct{}{}
	}
	go p.watchFile()
	return p, nil
}
func (p *Photo) close() {
	defer p.ctx.Done()
	defer p.watcher.Close()
}

// 添加监控路径
func (p *Photo) addPath(p1 string) error {
	if p.ctx.Err() != nil {
		return errors.New("context close")
	}
	p1, err := filepath.Abs(p1)
	if err != nil {
		return err
	}
	// 排除路径
	if _, ok := p.excludePaths[p1]; ok {
		return errors.New("exclude path")
	}
	// 排除文件夹名称
	if _, ok := p.excludeDirs[filepath.Base(p1)]; ok {
		return errors.New("exclude dirName")
	}
	if info, err := os.Lstat(p1); err != nil { //判断文件夹是否存在
		return err
	} else if !info.IsDir() {
		return errors.New("not found dir")
	}
	if err := p.watcher.Add(p1); err != nil {
		return err
	}
	// 监听深层目录
	return filepath.WalkDir(p1, func(p2 string, d fs.DirEntry, err error) error {
		if p.ctx.Err() != nil {
			return errors.New("context close")
		}
		if err != nil {
			return err
		}
		if d.IsDir() {
			return p.watcher.AddWith(p2)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return p.walkPhotos(p2, info)
	})
}
func (p *Photo) watchFile() {
	for {
		select {
		case err, ok := <-p.watcher.Errors:
			if !ok {
				return
			}
			if err != nil {
				log.Println("photos.watch:", err)
			}
		case event, ok := <-p.watcher.Events:
			if !ok {
				return
			}
			switch event.Op {
			//写入内容
			case fsnotify.Write:
			//模式
			case fsnotify.Chmod:
			//创建
			case fsnotify.Create:
				p.createFile(event.Name)
			// 删除文件
			case fsnotify.Remove:
				p.removeFile(event.Name)
			// 重命名文件
			case fsnotify.Rename:
				p.removeFile(event.Name)
			// 删除并重命名文件
			case fsnotify.Remove | fsnotify.Rename:
				p.removeFile(event.Name)
			}
		}
	}
}
func (p *Photo) createFile(p1 string) error {
	//判断文件是否存在
	info, err := os.Lstat(p1)
	if err != nil {
		return err
	}
	//判断是否文件夹
	if info.IsDir() {
		return p.addPath(p1)
	}
	// 处理里面的视频文件
	return p.walkPhotos(p1, info)
}
func (p *Photo) removeFile(p1 string) error {
	if p.ctx.Err() != nil {
		return errors.New("context close")
	}
	p.watcher.Remove(p1)
	return nil
}
func (p *Photo) genFileID(p1 string) string {
	return GenFileID(p.serialId, toBase64(p1))
}

// 搜索该目录下的所有照片视频文件
func (p *Photo) walkPhotos(p1 string, info fs.FileInfo) error {
	var err error
	if info == nil {
		if info, err = os.Lstat(p1); err != nil {
			return err
		}
	}
	if info.IsDir() {
		return nil
	}
	ext := strings.ToLower(strings.TrimPrefix(path.Ext(filepath.Base(p1)), "."))
	switch ext {
	// 照片
	case "jpg", "gif", "png", "webp", "jpeg", "heic", "heif":
		return p.addImageIndex(p1, info)
	// 视频
	case "mp4", "mkv", "mov", "wmv", "flv", "m3u8", "ts", "avi", "webm", "avchd":
		return p.addVideoIndex(p1, info)
	}
	return err
}

// 添加视频
func (p *Photo) addVideoIndex(p1 string, info fs.FileInfo) error {
	// 生成文件id
	fileid := p.genFileID(p1)
	// 判断是否已存在
	if ok := p.ctx.Search.Exist(fileid); ok {
		return nil
	}
	item := &SearchItem{
		SerialId:      p.serialId,
		Filename:      info.Name(),
		Path:          p1,
		Size:          info.Size(),
		FileType:      FileType_Video,
		LastDate:      info.ModTime().Format(time.DateTime),
		LastTimestamp: info.ModTime().Unix(),
	}
	// 生成封面
	GenVideoCover(p1, p.ctx.coverPath)
	if err := p.addSearch(fileid, item); err != nil {
		return err
	}
	return nil
}

// 添加照片
func (p *Photo) addImageIndex(p1 string, info fs.FileInfo) error {
	// 生成文件id
	fileid := p.genFileID(p1)
	// 判断是否已存在
	if ok := p.ctx.Search.Exist(fileid); ok {
		return nil
	}
	item := &SearchItem{
		SerialId:      p.serialId,
		Filename:      info.Name(),
		Path:          p1,
		Size:          info.Size(),
		FileType:      FileType_IMAGE,
		LastDate:      info.ModTime().Format(time.DateTime),
		LastTimestamp: info.ModTime().Unix(),
	}
	if rawExif, err := goexif.SearchFileAndExtractExif(p1); err == nil {
		if tags, _, err := goexif.GetFlatExifData(rawExif, &goexif.ScanOptions{}); err == nil {
			for i := 0; i < len(tags); i++ {
				switch strings.TrimSpace(tags[i].TagName) {
				case "ImageWidth":
					if v, ok := anyToInt64(tags[i].Value); ok {
						item.ExifWidth = v
					}
				case "ImageHeight":
					if v, ok := anyToInt64(tags[i].Value); ok {
						item.ExifHeight = v
					}
				case "Model":
					if v, ok := tags[i].Value.(string); ok {
						item.ExifModel = v
					}
				case "ImageLength":
					if v, ok := anyToInt64(tags[i].Value); ok {
						item.ExifLength = v
					}
				case "DateTime":
					if item.LastDate == "" {
						if v, ok := tags[i].Value.(string); ok {
							item.LastDate = v
						}
					}
				case "DateTimeOriginal":
					if v, ok := tags[i].Value.(string); ok {
						item.ExifOriginalDate = v
					}
				}
				// fmt.Printf("%v: %v\n", tags[i].TagName, tags[i].Value)
			}
		}
	}
	if img, err := imaging.Open(p1); err == nil {
		if p.ctx.Store != nil {
			hash, _ := duplo.CreateHash(img)
			p.ctx.Store.Add(p1, hash)
		}
		// 生成封面
		ImageToCover(img, filepath.Join(p.ctx.coverPath, GenCoverFilename(p1)), 220, 220)
	}
	return p.addSearch(fileid, item)
}

// 添加文件至索引
func (p *Photo) addSearch(fileid string, value *SearchItem) error {
	if p.ctx.Err() != nil {
		return errors.New("context close")
	}
	// if p.ctx.Search == nil {
	// 	return errors.New("uninitialized search")
	// }
	return p.ctx.Search.Add(map[string]*SearchItem{
		fileid: value,
	})
}

// 为文件添加标签
func (p *Photo) AddTag(p1 string, tag string) error {
	if p.ctx.Err() != nil {
		return errors.New("context close")
	}
	fileid := GenFileID(p1)
	datas, err := p.ctx.Search.Query(RequestSearch{
		Filters: map[string]interface{}{
			Index_SerialId: fileid,
		},
	})
	if err != nil {
		return err
	}
	items := make(map[string]*SearchItem)
	for i := 0; i < len(datas.Result); i++ {
		isExist := false
		tags := datas.Result[i].GetTags()
		for j := 0; j < len(tags); j++ {
			if tags[i] == tag {
				isExist = true
				break
			}
		}
		if !isExist {
			tags = append(tags, tag)
			datas.Result[i].Tags = tags
			items[fileid] = datas.Result[i]
		}
	}
	if len(items) > 0 {
		return p.ctx.Search.Add(items)
	}
	return nil
}
