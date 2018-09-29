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
		gradles := handler.LocalGradle()
		for k,v := range gradles {
			fmt.Println(k + " => " + v)
		}
		context.JSON(http.StatusOK, gradles)
	})

	r.Run(":8090")
}