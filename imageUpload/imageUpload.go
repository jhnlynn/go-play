package imageUpload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
	"log"
	"net/http"
	"os"
	"strings"
)

// UPLOAD THUMBNAIL upload directories
const UPLOAD = "./uploads/"
const THUMBNAIL = "./uploads/thumbs/"

// UploadImageAndThumb
// upload image received from request,
// and create & store thumbnail generated based on this image
///*
func UploadImageAndThumb(c *gin.Context) {
	log.Println("start to upload image...")

	allowCORS(c)

	// upload the image received from request
	img := saveImage(c)

	generateAndSaveThumbnail(c, img, 300, 300)
}

// general error responding
func processErr(c *gin.Context, err error, state string) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf(state, err),
		})
		return
	}
}

// allow CORS
func allowCORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

}

// saveImage:
// save image, and process error
func saveImageGeneral(c *gin.Context, img *imageupload.Image, dir string) {
	log.Println("create the directory")

	// if not exists, create the 'uploads' directory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": fmt.Sprintf("image saving failed, %v", err),
		})
		panic(err)
	}

	log.Println("Create file and save the image")

	// Create file
	dst, err1 := os.Create(dir + img.Filename)
	defer dst.Close()

	err2 := img.Save(dir + img.Filename)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error": fmt.Sprintf("image saving failed"),
		})
		return
	}
}

// saveImage:
// save image from request
func saveImage(c *gin.Context) *imageupload.Image {
	img, err := imageupload.Process(c.Request, "photo")

	processErr(c, err, "image multi-part uploading error: %v")

	fileInfo := strings.Split(img.Filename, ".")
	ext := fileInfo[1]

	// if it is not an image, response format error to web
	if ext != "jpeg" && ext != "png" && ext != "gif" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("file format is not correct, please uplaod a jpeg/png/gif image"),
		})
		return nil
	}

	log.Println("Save image itself")
	saveImageGeneral(c, img, UPLOAD)

	return img
}

// generateAndSaveThumbnail:
// generate thumbnail using 'imageupload' package,
// and save it
func generateAndSaveThumbnail(c *gin.Context, img *imageupload.Image, width int, height int) {
	log.Println("generate thumbnail and save it")

	thumb, err := imageupload.ThumbnailPNG(img, width, height)

	processErr(c, err, "image thumbnail producing error: %v")

	imgInfo := strings.Split(img.Filename, ".")

	thumb.Filename = imgInfo[0] + "-thumbnail" + ".png"

	saveImageGeneral(c, thumb, THUMBNAIL)
}
