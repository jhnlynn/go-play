package redisOp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strings"
	"time"
)

func RedisNewClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:	  "localhost:6379",
		Password: "", // no password set
		DB:		  0,  // use default DB
	})

	return rdb
}

// GetRedis return if the key has been found in redis
func GetRedis(c *gin.Context, blogId string) bool {

	strCmd := RedisNewClient().Get(c, blogId)
	if err := strCmd.Err(); err != nil {
		errBody := fmt.Sprintf("%v", err)
		if info := strings.Split(errBody, ":")[1]; info == " nil" {
			return false
		} else {
			c.JSON(http.StatusInternalServerError, gin.H {
				"error": err,
			})
			return false
		}
	}

	return true
}

func SetRedis(c *gin.Context, blogId string) bool {

	ipAddr, _ := c.RemoteIP()

	// [ipAddr]: ipAddr, expiration: 20 s
	boolCmd := RedisNewClient().Set(c, blogId, ipAddr.String(), 20*time.Second)
	fmt.Println("boolCmd: ", boolCmd)
	if err := boolCmd.Err(); err != nil {
		fmt.Println("redis HSet", err)
		c.JSON(http.StatusBadRequest, gin.H {
			"error": err,
		})
		return false
	}

	return true
}
