package handler

import (
	"container/list"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"os/user"
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
	LocalGradle() map[string]string
	GradleCacheList() (map[string]list.List, error)
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
func (osHandler MacOSHandler) LocalGradle() map[string]string {

	gradles := make(map[string]string)
	// 查找android studio的
	as := "/Applications/Android Studio.app/Contents/gradle"
	gradlesInAS := gradleInDir(as)
	for k,v := range gradlesInAS {
		gradles[k] = v
	}
	// 查找根目录下的gradle
	username, err := user.Current()
	var rootDir = ""
	if err == nil {
		rootDir = "/Users/" + username.Name
		gradlesInRoot := gradleInDir(rootDir)
		for k,v := range gradlesInRoot {
			gradles[k] = v
		}
	}
	return gradles
}

func (osHandler MacOSHandler) GradleCacheList() (map[string]list.List, error) {

	gradleVersionsMap := make(map[string]list.List)

	// 查找根目录下的gradle缓存， ~/Users/xx/.gradle/caches/jars-3
	username, err := user.Current()

	if err != nil {
		return nil, err
	}

	gradleCacheDir := "/Users/" + username.Username + "/.gradle/caches/jars-3"
	gradleCacheDirInfo, err := os.Stat(gradleCacheDir)

	if err != nil {
		return nil, err
	}

	if gradleCacheDirInfo.IsDir() {
		//  确定是文件夹
		dir, err := ioutil.ReadDir(gradleCacheDir)
		if err != nil {
			return nil, err
		}

		for _, dirItem := range dir {
			//  解析库版本

			jarVersionMap, err := ParseGradleJars(gradleCacheDir, dirItem.Name())

			if err != nil {
				return nil, err
			}

			for k, v := range jarVersionMap {
				gradleVersionsMap[k] = v
			}

		}
	}

	return gradleVersionsMap, nil
}

// 查找某个父目录下的gradle目录
func gradleInDir(parent string) map[string]string {

	var (
		fileInfo       os.FileInfo
		err            error
		gradlesMap = make(map[string]string)
	)

	fileInfo, err = os.Stat(parent)

	if err != nil {
		fmt.Println("error! ", err)
	} else {
		if fileInfo.IsDir() {
			// 是文件夹
			dir, err := ioutil.ReadDir(parent)
			if err != nil {
				println("Android Studio gradle dir error!")
			} else {
				for _, fi := range dir {
					fileName := fi.Name()
					if dirIsGradle(parent, fileName) {
						// 是 gradle 文件夹
						gradleVersion := utils.GradleVersion(fileName)
						// 完整的路径名
						gradleDirCompleteName := parent + string(filepath.Separator) + fileName
						gradlesMap[gradleDirCompleteName] = gradleVersion
						println(fi.Name() + "是gradle文件夹, 版本是 " + gradleVersion)
					}
				}
			}
		}
	}
	//for k,v := range gradlesMap{
	//	println(k + " => " + v)
	//}
	return gradlesMap
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



// 根据gradle 缓存的jar包的名称和父目录去解析gradle库的版本

func ParseGradleJars(parent string, jarDirName string) (map[string]list.List, error) {
	jarVersionMap := make(map[string]list.List)
	finalJarDirName := parent + string(filepath.Separator) + jarDirName
	dir, err := os.Stat(finalJarDirName)

	if err != nil {
		return nil ,err
	}

	if dir.IsDir() {
		dir, err := ioutil.ReadDir(finalJarDirName)

		if err  != nil {

			return nil, err

		}

		// 遍历dir文件夹，处理里面每一个jar包

		for _, dirItem := range dir {
			// 处理每一个jar包

			jarName := dirItem.Name()
			jarVersionMap[jarName] = handleJar(finalJarDirName, jarName)
		}

	}

	return jarVersionMap, nil
}



// 解析具体的jar 包

func handleJar(parent, jarName string) (list.List) {
	println(parent+ " => " + jarName)
	return *list.New()
}