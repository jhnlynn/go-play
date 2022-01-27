package main

import (
	"github.com/joho/godotenv"
	"go-play/routers"
	"log"
	"os"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
		os.Exit(1)
	}
}

func main() {
	LoadEnv()

	//awsAccessKeyID := GetEnvWithKey("AWS_ACCESS_KEY_ID")
	//fmt.Println("My access key ID is ", awsAccessKeyID)

	router := routers.Routers()
	err := router.SetTrustedProxies([]string{"192.168.0.2"})
	if err != nil {
		panic(err)
		return
	}

	_ = router.Run()
}
