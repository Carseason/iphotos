package iphotos

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

func GenCoverFilename(inputFilePath string) string {
	return toBase64(inputFilePath) + ".png"
}

// 生成照片略缩图
// 保存路径，生成的文件
func GenCover(outPath, inputFilePath string, width, height int) (string, error) {
	var err error
	outPath, err = filepath.Abs(outPath)
	if err != nil {
		return "", err
	}
	inputFilePath, err = filepath.Abs(inputFilePath)
	if err != nil {
		return "", err
	}
	abspath := filepath.Join(outPath, GenCoverFilename(inputFilePath))
	// // 判断文件是否已存在
	if _, err := os.Stat(abspath); err == nil {
		return abspath, nil
	}
	r, err := imaging.Open(inputFilePath)
	if err != nil {
		return "", err
	}
	return ImageToCover(r, abspath, width, height)
}
func ImageToCover(r image.Image, abspath string, width, height int) (string, error) {
	img := imaging.Thumbnail(r, width, height, imaging.Lanczos)
	err := imaging.Save(img, abspath)
	// 如果是目录未找到
	if err != nil && errors.Is(err, os.ErrNotExist) {
		outPath, _ := filepath.Split(abspath)
		if err = mkdirAll(outPath); err != nil {
			return "", err
		}
		err = imaging.Save(img, abspath)
	}
	return abspath, err
}

// 流
func WriterCover(w io.Writer, inputFilePath string, width, height int) error {
	format, err := imaging.FormatFromFilename(inputFilePath)
	if err != nil {
		return err
	}
	r, err := imaging.Open(inputFilePath)
	if err != nil {
		return err
	}
	img := imaging.Thumbnail(r, width, height, imaging.Lanczos)
	err = imaging.Encode(w, img, format)
	return err
}
func GenVideoCover(inputPath, outPath string, ns ...int) (string, error) {
	filename := GenCoverFilename(inputPath)
	abspath := filepath.Join(outPath, filename)
	// // 判断文件是否已存在
	if _, err := os.Stat(abspath); err == nil {
		return abspath, nil
	}
	var args []string
	if len(ns) == 0 { //关键帧
		args = []string{
			"-skip_frame",
			"nokey",
			"-i",
			inputPath,
			"-vsync",
			"0",
			"-f",
			"image2",
			"-vcodec",
			"mjpeg",
			"-vframes",
			"1",
			"-y",
		}
	} else { // 第几秒
		frameStr := fmt.Sprintf("%d", ns[0])
		args = []string{
			// 把时间放在文件前，避免等待帧,速度更快,但是可能时间不存在而失败
			"-ss",
			frameStr,
			"-i",
			inputPath,
			"-f",
			"image2",
			"-vframes",
			"1",
			"-y",
		}
	}
	// 30s 后不管结果都kill掉
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	args = append(args, []string{
		abspath,
	}...)
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", errors.New(stderr.String())
	}
	return abspath, nil
}
