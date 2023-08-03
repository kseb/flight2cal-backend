package utils

import "os"

func AirlabsToken() string {
	return os.Getenv("AIRLABS_TOKEN")
}

func AllowOriginHost() string {
	return os.Getenv("ALLOW_ORIGIN_HOST")
}
