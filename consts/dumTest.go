package consts

import "go-play/common/getEnv"

func GetMongoAPI() string {
	return getEnv.EnvWithKey("MONGO_URI")
}