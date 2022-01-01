package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var filepath string

func UploadImage(c *gin.Context) {
	log.Println("start to upload image...")

	sess := c.MustGet("sess").(*session.Session)

	uploader := s3manager.NewUploader(sess)
	MyBucket := GetEnvWithKey("BUCKET_NAME")

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		fmt.Println("error from c.Request.FormFile: \n", err)
	}

	filename := uuid.New().String() + "." +  strings.Split(header.Filename, ".")[1]

	log.Println("file name: \n", filename)

	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput {
		Bucket: aws.String(MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"error":    "Failed to upload file",
			"uploader": up,
		})
		fmt.Println("failed to upload file")
		fmt.Println("here's the error: ", err)
		return
	}
	filepath = "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/" + filename
	c.JSON(http.StatusOK, gin.H{
		"filepath":    filepath,
	})
}

func ConnectAws() *session.Session {
	AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
	MyRegion = GetEnvWithKey("AWS_REGION")
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}

//GetEnvWithKey : get env value
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func main() {
	LoadEnv()

	//awsAccessKeyID := GetEnvWithKey("AWS_ACCESS_KEY_ID")
	//fmt.Println("My access key ID is ", awsAccessKeyID)

	sess := ConnectAws()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("sess", sess)
		c.Next()
		CORSMiddleware()
	})

	router.POST("/upload", UploadImage)

	_ = router.Run(":4000")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}