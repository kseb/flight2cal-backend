package utils

import (
	"github.com/gin-gonic/gin"
)

func AddAccessControlAllowOriginIfSet(context *gin.Context) {
	if AllowOriginHost() != "" {
		context.Header("Access-Control-Allow-Origin", AllowOriginHost())
	}
}
