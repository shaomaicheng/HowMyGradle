package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/myGradle/src/handler"
	"runtime"
)



func main() {

	osHandlerManager := handler.GetInstance()

	osHandlerManager.RegisterOS("darwin",  handler.MacOSHandler{})


	r := gin.Default()
	osHandlerManager.Dispatch(r)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprint(runtime.GOOS))
		handler := c.MustGet(handler.OS_HANDLER).(handler.OSHandler)
		handler.Root()
	})

	r.GET("/localgradle", func(context *gin.Context) {
		handler := context.MustGet(handler.OS_HANDLER).(handler.OSHandler)
		handler.LocalGradle()
		context.String(http.StatusOK, fmt.Sprint("local_gradle"))
	})

	r.Run(":8090")
}