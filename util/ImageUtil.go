package util

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
)

// Resize 图片缩放
/*
 * lanczos- 用于产生清晰结果的摄影图像的高质量重采样过滤器。
 * catmullRom- 用于产生动画的高质量重采样过滤器。
 * linear- 用于产生平滑结果的快速重采样过滤器。
 * box- 用于产生粗糙结果的快速重采样过滤器。
 * nearestNeighbor- 用于产生粗糙结果的最快重采样过滤器。
 * gaussian- 用于产生平滑结果的重采样过滤器。
 * mitchellNetravali- 用于产生动画的高质量重采样过滤器。
 * hermite- 用于产生平滑结果的重采样过滤器。
 * bSpline- 用于产生平滑结果的重采样过滤器。
 * hann- 用于产生平滑结果的重采样过滤器。
 * hamming- 用于产生平滑结果的重采样过滤器。
 * blackman- 用于产生平滑结果的重采样过滤器。
 */
func Resize(file *multipart.FileHeader, width, height int, filter string) *multipart.FileHeader {
	var resizeFilter = map[string]imaging.ResampleFilter{
		"lanczos":           imaging.Lanczos,
		"catmullRom":        imaging.CatmullRom,
		"linear":            imaging.Linear,
		"box":               imaging.Box,
		"nearestNeighbor":   imaging.NearestNeighbor,
		"gaussian":          imaging.Gaussian,
		"mitchellNetravali": imaging.MitchellNetravali,
		"hermite":           imaging.Hermite,
		"bspline":           imaging.BSpline,
		"hann":              imaging.Hann,
		"hamming":           imaging.Hamming,
		"blackman":          imaging.Blackman,
	}

	img, imgType, filename := FileToImage(file)

	img = imaging.Resize(img, width, height, resizeFilter[filter])

	return ImageToFile(img, imgType, filename)

}

func FileToImage(header *multipart.FileHeader) (img image.Image, imgType string, filename string) {
	filename = header.Filename
	file, _ := header.Open()
	bufFile := bufio.NewReader(file)
	img, imgType, _ = image.Decode(bufFile)

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("转换错误")
			return
		}
	}(file)

	return img, imgType, filename
}

func ImageToFile(img image.Image, imgType, filename string) *multipart.FileHeader {
	bufFile := &bytes.Buffer{}

	if imgType == "jpg" {
		// 保存裁剪的图片
		_ = jpeg.Encode(bufFile, img, nil)
	} else if imgType == "png" {
		// 保存裁剪的图片
		_ = png.Encode(bufFile, img)
	}
	newFile, _ := os.Create(filename)
	fi, _ := newFile.Stat()

	_, err := io.Copy(newFile, bufFile)
	if err != nil {
		fmt.Printf("转换错误")
		return nil
	}

	return &multipart.FileHeader{
		Filename: fi.Name(),
		Size:     fi.Size(),
	}
}
