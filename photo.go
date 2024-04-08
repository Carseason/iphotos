package iphotos

import (
	"errors"
	"io/fs"
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
	// 从 photos 传来的
	ctx *Context
	mx  *sync.RWMutex
	//是否正在配置中
	indexing bool
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

// 关闭实例
func (p *Photo) close() {
	p.ctx.cancel()
	p.mx.Lock()
	defer p.mx.Unlock()
	p.watcher.Close()
	p.indexing = false
}

// 索引状态,不是强业务,不加锁
func (p *Photo) Indexing() bool {
	return p.indexing
}

// 添加路径, 创建了实例后使用该方法添加路径
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
	// 避免堵塞
	go p.startIndex(p1)
	return nil
}

// 开始索引,尽量使用异步开始,并使用 indexing 来判断是否完成
func (p *Photo) startIndex(p1 string) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.indexing = true
	p.addPath(p1)
	p.indexing = false
}

// 添加照片库路径
// 必须是文件夹,并且监控该路径
func (p *Photo) addPath(p1 string) error {
	if p.ctx.Err() != nil {
		return ErrContextClose
	}
	// 排除路径
	if _, ok := p.excludePaths[p1]; ok {
		return errors.New("exclude path")
	}
	// 排除文件夹名称
	if _, ok := p.excludeDirs[filepath.Base(p1)]; ok {
		return errors.New("exclude dirName")
	}
	// 监控文件夹至
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
		// 先处理照片在处理文件夹
		if n := len(files); n > 0 {
			for i := 0; i < n; i++ {
				p2 := filepath.Join(p1, files[i].Name())
				if files[i].IsDir() {
					dirs = append(dirs, p2)
				} else {
					p.walkPhotos(p2)
				}
			}
		}
		if err != nil {
			break
		}
	}
	// 已经处理完照片文件里,手动关闭当前文件
	// 再开始处理文件夹，避免长期占着文件句柄
	f.Close()
	// 处理文件夹,一层一层处理
	for i := 0; i < len(dirs); i++ {
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
				// 已经关闭了
				if errors.Is(err, fsnotify.ErrClosed) {
					return
				}
				log.Println("photos.watch:", err)
			}
		case event, ok := <-p.watcher.Events:
			if !ok {
				return
			}
			switch event.Op {
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

// 创建文件
func (p *Photo) createFile(p1 string) error {
	//判断文件是否存在
	info, err := os.Lstat(p1)
	if err != nil {
		return err
	}
	//判断是否文件夹,把文件夹添加至监控
	if info.IsDir() {
		return p.addPath(p1)
	}
	// 如果不是文件夹,则处理该文件
	return p.walkPhotos(p1)
}

// 删除文件
func (p *Photo) removeFile(p1 string) error {
	if p.ctx.Err() != nil {
		return ErrContextClose
	}
	// 取消监控
	p.watcher.Remove(p1)
	fileid, err := GenFileID(p.serialId, p1)
	if err != nil {
		return err
	}
	// 从索引里删除该文件
	p.ctx.Delete(fileid)
	return nil
}

// 处理照片视频文件
func (p *Photo) walkPhotos(p1 string) error {
	if p.ctx.Err() != nil {
		return ErrContextClose
	}
	ext := strings.ToLower(strings.TrimPrefix(path.Ext(filepath.Base(p1)), "."))
	switch ext {
	// 照片
	case "jpg", "gif", "png", "webp", "jpeg", "heic", "heif":
		return p.addImageIndex(p1, ext)
	// 视频
	case "mp4", "mkv", "mov", "wmv", "flv", "avi", "webm", "avchd":
		return p.addVideoIndex(p1, ext)
	}
	return nil
}

// 检查文件是否已存在
func (p *Photo) checkFile(p1 string) (string, fs.FileInfo, error) {
	// 生成文件id
	fileid, err := GenFileID(p.serialId, p1)
	if err != nil {
		return "", nil, err
	}
	// 判断是否已存在
	if ok := p.ctx.Exist(fileid); ok {
		return fileid, nil, nil
	}
	// 获取文件信息
	info, err := os.Stat(p1)
	if err != nil {
		return "", nil, err
	}
	return fileid, info, nil
}

// 添加视频
func (p *Photo) addVideoIndex(p1, ext string) error {
	fileid, info, err := p.checkFile(p1)
	if err != nil {
		return err
	}
	if info == nil {
		return nil
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
	}
	// 添加至搜索引擎
	return p.ctx.Add(map[string]*SearchItem{
		fileid: item,
	})
}

// 添加照片
func (p *Photo) addImageIndex(p1, ext string) error {
	fileid, info, err := p.checkFile(p1)
	if err != nil {
		return err
	}
	if info == nil {
		return nil
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
	}
	// 处理exif
	if rawExif, err := GetImageExif(p1); err == nil {
		item.ExifHeight = rawExif.ExifHeight
		item.ExifWidth = rawExif.ExifWidth
		item.ExifModel = rawExif.ExifModel
		item.ExifOriginalDate = rawExif.ExifOriginalDate
		item.ExifMake = rawExif.ExifMake
	}
	// 添加至搜索引擎
	return p.ctx.Add(map[string]*SearchItem{
		fileid: item,
	})
}
