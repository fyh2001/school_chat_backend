package controllers

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"os"
	"path"
	Results "schoolChat/app/result"
)

//var uploadPath = "/Users/huangye/Downloads/upload/"

var uploadPath = "/www/wwwroot/43.139.54.138/schoolWall/upload/"

// UploadHandler 上传文件
func UploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Printf(err.Error())
		c.JSON(200, Results.Err.Fail("上传失败"))
		return
	}
	files := form.File["file"]

	var successFiles []string

	for _, file := range files {
		UUID := uuid.NewV4()
		fileSuffixName := path.Ext(file.Filename)
		filename := UUID.String() + fileSuffixName

		if file != nil {
			// 保存上传的文件到临时路径
			tempFilePath := uploadPath + "temp_" + filename
			if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
				fmt.Printf(err.Error())
				c.JSON(200, Results.Err.Fail("上传失败"))
				return
			}

			// 打开临时文件
			tempFile, err := imaging.Open(tempFilePath)
			if err != nil {
				fmt.Printf(err.Error())
				c.JSON(200, Results.Err.Fail("上传失败"))
				return
			}

			// 进行图片修改操作
			//resizedImage := imaging.Resize(tempFile, 800, 800, imaging.Lanczos)
			// 将图片转换为JPEG格式
			// 创建一个新的文件
			outputFile, err := os.Create(uploadPath + filename)
			if err != nil {
				fmt.Printf(err.Error())
				c.JSON(200, Results.Err.Fail("上传失败"))
				return
			}
			defer outputFile.Close()

			// 将修改后的图片写入文件
			if err := imaging.Encode(outputFile, tempFile, imaging.JPEG, imaging.JPEGQuality(30)); err != nil {
				fmt.Printf(err.Error())
				c.JSON(200, Results.Err.Fail("上传失败"))
				return
			}

			//// 保存修改后的图片到最终路径
			//if err := imaging.Save(resizedImage, uploadPath+filename); err != nil {
			//	fmt.Printf(err.Error())
			//	c.JSON(200, Results.Err.Fail("上传失败"))
			//	return
			//}

			// 删除临时文件
			if err := os.Remove(tempFilePath); err != nil {
				fmt.Printf(err.Error())
				// 可以选择忽略临时文件删除错误
			}

			//if err := c.SaveUploadedFile(tempFileHeader, uploadPath+filename); err != nil {
			//	fmt.Printf(err.Error())
			//	c.JSON(200, Results.Err.Fail("上传失败"))
			//	return
			//}
		} else {
			c.JSON(200, Results.Err.Fail("上传失败"))
			return
		}

		successFiles = append(successFiles, filename)
	}

	c.JSON(200, successFiles)
}

// DownloadHandler 下载文件
func DownloadHandler(c *gin.Context) {
	name := c.Query("filename")

	if len([]byte(name)) == 0 {
		c.JSON(200, Results.Err.Fail("文件名不能为空"))
	}

	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+name)
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(uploadPath + name)
}

// DownloadAvatarHandler 下载头像
func DownloadAvatarHandler(c *gin.Context) {
	name := c.Query("filename")

	if len([]byte(name)) == 0 {
		c.JSON(200, Results.Err.Fail("文件名不能为空"))
	}

	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+name)
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Transfer-Encoding", "binary")
	c.File(uploadPath + "avatar/" + name)
}
