package iphotos

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Photo struct {
	watcher  *fsnotify.Watcher
	serialId string
	// 排除的文件夹名称,在里面的文件夹和子文件夹则排除
	excludeDirs map[string]struct{}
	// 排除的路径
	excludePaths map[string]struct{}
	ctx          *Context
	indexing     bool //是否正在配置中
	paths        []string
	mx           *sync.RWMutex
	total        int64 //统计数量
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
		mx:           &sync.RWMutex{},
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
	p.ctx.Done()
	p.mx.Lock()
	defer p.mx.Unlock()
	p.total = 0
	p.watcher.Close()
}
func (p *Photo) Indexing() bool {
	return p.indexing
}
func (p *Photo) Total() int64 {
	return p.total
}
func (p *Photo) AddPATH(p1 string) error {
	p.mx.Lock()
	defer p.mx.Unlock()
	p1, err := filepath.Abs(p1)
	if err != nil {
		return err
	}
	if info, err := os.Lstat(p1); err != nil { //判断文件夹是否存在
		return err
	} else if !info.IsDir() {
		return errors.New("not found dir")
	}
	//
	for i := 0; i < len(p.paths); i++ {
		if p.paths[i] == p1 {
			return nil
		}
	}
	p.paths = append(p.paths, p1)
	// 避免堵塞
	go p.startIndex(p1)
	return nil
}
func (p *Photo) startIndex(p1 string) {
	p.mx.Lock()
	defer func() {
		p.indexing = false
		p.mx.Unlock()
	}()
	p.indexing = true
	p.addPath(p1)
}

// 重新开始索引
func (p *Photo) Restart() error {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.indexing = true
	defer func() {
		p.indexing = false
	}()
	var errs error
	for i := 0; i < len(p.paths); i++ {
		if err := p.addPath(p.paths[i]); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

// 添加照片库路径
// 必须是文件夹,并且监控该路径
func (p *Photo) addPath(p1 string) error {
	if p.ctx.Err() != nil {
		return errors.New("context close")
	}
	// 排除路径
	if _, ok := p.excludePaths[p1]; ok {
		return errors.New("exclude path")
	}
	// 排除文件夹名称
	if _, ok := p.excludeDirs[filepath.Base(p1)]; ok {
		return errors.New("exclude dirName")
	}
	if err := p.watcher.Add(p1); err != nil {
		return err
	}
	f, err := os.Open(p1)
	if err != nil {
		return err
	}
	defer f.Close()
	dirs := []string{}
	for {
		// 每次只读10条
		files, err := f.Readdir(10)
		if err != nil {
			break
		}
		n := len(files)
		// 先处理照片在处理文件夹
		for i := 0; i < n; i++ {
			p2 := filepath.Join(p1, files[i].Name())
			if files[i].IsDir() {
				dirs = append(dirs, p2)
			} else {
				p.walkPhotos(p2)
			}
		}
	}
	// 已经处理完照片文件里,手动关闭当前文件
	// 再开始处理文件夹
	f.Close()
	// 处理文件夹
	ds := len(dirs)
	for i := 0; i < ds; i++ {
		p.addPath(dirs[i])
	}
	return nil
}

// 监控文件
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

// 生成文件ID
func (p *Photo) genFileID(p1 string) (string, error) {
	return p.ctx.GenFileID(p.serialId, p1)
}

// 创建文件
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
	return p.walkPhotos(p1)
}

// 删除文件
func (p *Photo) removeFile(p1 string) error {
	if p.ctx.Err() != nil {
		return errors.New("context close")
	}
	fileid, err := p.genFileID(p1)
	if err != nil {
		return err
	}
	// 取消监控
	p.watcher.Remove(p1)
	// 从索引里删除
	p.ctx.Delete(fileid)
	return nil
}

// 处理照片视频文件
func (p *Photo) walkPhotos(p1 string) error {
	if p.ctx.Err() != nil {
		return errors.New("context close")
	}
	ext := strings.ToLower(strings.TrimPrefix(path.Ext(filepath.Base(p1)), "."))
	switch ext {
	// 照片
	case "jpg", "gif", "png", "webp", "jpeg", "heic", "heif":
		return p.addImageIndex(p1, ext)
	// 视频
	case "mp4", "mkv", "mov", "wmv", "flv", "m3u8", "ts", "avi", "webm", "avchd":
		return p.addVideoIndex(p1, ext)
	}
	return nil
}

// 添加视频
func (p *Photo) addVideoIndex(p1, ext string) error {
	// 生成文件id
	fileid, err := p.genFileID(p1)
	if err != nil {
		return err
	}
	// 判断是否已存在
	if ok := p.ctx.Exist(fileid); ok {
		return nil
	}
	info, err := os.Stat(p1)
	if err != nil {
		return err
	}
	item := &SearchItem{
		SerialId:      p.serialId,
		Filename:      info.Name(),
		Path:          p1,
		Size:          strconv.FormatInt(info.Size(), 10),
		FileType:      FileType_Video,
		FileExt:       ext,
		LastDate:      info.ModTime().Format(time.DateTime),
		LastTimestamp: strconv.FormatInt(info.ModTime().Unix(), 10),
		Public:        Public_PUBLIC,
	}
	p.total++
	return p.ctx.Add(map[string]*SearchItem{
		fileid: item,
	})
}

// 添加照片
func (p *Photo) addImageIndex(p1, ext string) error {
	// 生成文件id
	fileid, err := p.genFileID(p1)
	if err != nil {
		return err
	}
	// 判断是否已存在
	if ok := p.ctx.Exist(fileid); ok {
		return nil
	}
	info, err := os.Stat(p1)
	if err != nil {
		return err
	}
	item := &SearchItem{
		SerialId:      p.serialId,
		Filename:      info.Name(),
		Path:          p1,
		Size:          strconv.FormatInt(info.Size(), 10),
		FileType:      FileType_IMAGE,
		FileExt:       ext,
		LastDate:      info.ModTime().Format(time.DateTime),
		LastTimestamp: strconv.FormatInt(info.ModTime().Unix(), 10),
		Public:        Public_PUBLIC,
	}
	// 处理exif
	if rawExif, err := GetImageExif(p1); err == nil {
		item.ExifHeight = rawExif.ExifHeight
		item.ExifWidth = rawExif.ExifWidth
		item.ExifModel = rawExif.ExifModel
		item.ExifOriginalDate = rawExif.ExifOriginalDate
		item.ExifMake = rawExif.ExifMake
	}
	p.total++
	return p.ctx.Add(map[string]*SearchItem{
		fileid: item,
	})
}
