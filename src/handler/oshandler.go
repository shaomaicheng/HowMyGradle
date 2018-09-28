package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path/filepath"
	"project/myGradle/src/utils"
	"regexp"
	"sync"
)

const MACOSX  = "darwin"
const OS_HANDLER  = "os_handler"


type OSHandlerManager struct {
	Handlers map[string]OSHandler
}

var manager * OSHandlerManager

func GetInstance() * OSHandlerManager {
	new(sync.Once).Do(func() {
		// just do once
		manager = new(OSHandlerManager)
		manager.Handlers = make(map[string]OSHandler)
	})
	return manager
}

func (osHandlerManager *OSHandlerManager) RegisterOS(name string, value OSHandler) {
	handlers := osHandlerManager.Handlers
	if handlers != nil {
		handlers[name] = value
	}
}


func (OSHandlerManager *OSHandlerManager) Dispatch(engine *gin.Engine) {
	handlers := OSHandlerManager.Handlers
	if handlers != nil {
		macHandler := handlers[MACOSX]
		engine.Use(func(context *gin.Context) {
			context.Set(OS_HANDLER, macHandler)
		})
	}
}

type OSHandler interface {
	Root()
	LocalGradle()
}


type MacOSHandler struct {
	OSHandler
}

func (osHandler MacOSHandler) Root() {
	fmt.Println("mac os")
}

/**
查找mac osx 本地的gradle
 */
func (osHandler MacOSHandler) LocalGradle() {
	// 查找android studio的
	as := "/Applications/Android Studio.app/Contents/gradle"

	// 根目录的
	//root := "~"


	var (
		fileInfo       os.FileInfo
		err            error
		gradlesMap = make(map[string]string)
	)

	fileInfo, err = os.Stat(as)

	if err != nil {
		fmt.Println("error! ", err)
	} else {
		if fileInfo.IsDir() {
			// 是文件夹
			dir, err := ioutil.ReadDir(as)
			if err != nil {
				println("Android Studio gradle dir error!")
			} else {
				for _, fi := range dir {
					fileName := fi.Name()
					if dirIsGradle(as, fileName) {
						// 是 gradle 文件夹
						gradleVersion := utils.GradleVersion(fileName)
						// 完整的路径名
						gradleDirCompleteName := as + string(filepath.Separator) + fileName
						gradlesMap[gradleDirCompleteName] = gradleVersion
						println(fi.Name() + "是gradle文件夹, 版本是 " + gradleVersion)
					}
				}
			}
		}
	}


	for k,v := range gradlesMap{
		println(k + " => " + v)
	}
}

func dirIsGradle(parent, dirname string) bool {
	dir, err := os.Stat(parent + string(filepath.Separator) + dirname)
	if err != nil {
		return false
	} else {
		return dir.IsDir() && isGradle(dirname)
	}
}

func isGradle(dirname string) bool {
	regex := "gradle-[1-9]+.[1-9]+"
	match, err := regexp.MatchString(regex, dirname)
	if err != nil {
		return false
	} else {
		return match
	}

}

