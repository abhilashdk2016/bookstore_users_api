package app

import (
	"github.com/abhilashdk2016/bookstore_utils_go/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApp() {
	mapUrls()
	logger.Info("About to start application!!!")
	router.Run("localhost:8081")
}
