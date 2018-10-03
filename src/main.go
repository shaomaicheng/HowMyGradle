package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/myGradle/src/oshandler"
	"runtime"
	"project/myGradle/src/model"
)



func main() {

	osHandlerManager := oshandler.GetInstance()

	osHandlerManager.RegisterOS("darwin",  oshandler.MacOSHandler{})


	r := gin.Default()
	osHandlerManager.Dispatch(r)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprint(runtime.GOOS))
		handler := c.MustGet(oshandler.OS_HANDLER).(oshandler.OSHandler)
		handler.Root()
	})

	r.GET("/localgradle", func(context *gin.Context) {
		handler := context.MustGet(oshandler.OS_HANDLER).(oshandler.OSHandler)
		gradles := handler.LocalGradle()
		for k,v := range gradles {
			fmt.Println(k + " => " + v)
		}
		context.JSON(http.StatusOK, gradles)
	})


	r.GET("/cachelist", func(context *gin.Context) {
		handler := context.MustGet(oshandler.OS_HANDLER).(oshandler.OSHandler)
		cacheList, err := handler.GradleCacheList()
		if err != nil {
			context.JSON(http.StatusExpectationFailed, "获取gradle jar缓存失败")
		} else {
			println(cacheList.Len())
			for iter := cacheList.Front(); iter != nil; iter = iter.Next() {
				item := iter.Value.(model.JarCache)
				println(item.Name)
			}
			context.JSON(http.StatusOK, cacheList)
		}
	})

	r.Run(":8090")
}