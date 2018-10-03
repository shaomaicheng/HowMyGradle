package oshandler

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
	"strings"
	"strconv"
	"project/myGradle/src/model"
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
	GradleCacheList() (*list.List, error)
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

// map[stirng]list.List  string=>库名称 List=>版本号数组
func (osHandler MacOSHandler) GradleCacheList() (*list.List, error) {

	gradleVersionsList := new(list.List)

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

			for k,v := range jarVersionMap {
				jarCacheItem := model.JarCache{k, v}
				gradleVersionsList.PushBack(jarCacheItem)
			}

		}
	}
	return gradleVersionsList, nil
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

		if err != nil {

			return nil, err

		}

		// 遍历dir文件夹，处理里面每一个jar包目录

		jarVersionList := new(list.List)
		jarName := ""
		jarWithoutVersion := ""
		jarVersion := ""

		for _, dirItem := range dir {
			// 处理每一个jar包目录

			jarName = dirItem.Name()
			jar := handleJar(jarName)
			jarWithoutVersion, jarVersion = getJarWithoutVersion(jar)
			if jar != "" {
				jarVersionList.PushBack(jarVersion)
			}
		}

		// 此处map的key 带有版本号和.jar， value的list里面带有版本号，需要再处理一次
		if jarWithoutVersion != "" {
			jarVersionMap[jarWithoutVersion] = *jarVersionList
		}

	}

	return jarVersionMap, nil
}



// 解析具体的jar 包

func handleJar(jarName string) string {
	isJar, jar := isJar(jarName)

	if !isJar {
		return ""
	}

	return jar
}

func isJar(jarName string) (bool,string) {
	jarRegex := "(.*).jar"
	isJar, err := regexp.MatchString(jarRegex, jarName)
	if err != nil {
		return false, ""
	}

	if !isJar {
		return false, ""
	}


	return isJar, jarName[0:len(jarName)-4]
}

// jar去掉版本号，返回jar自己的名称和版本号
func getJarWithoutVersion(jar string) (string,string) {
	jarVersionSplit := strings.Split(jar, "-")
	len := len(jarVersionSplit)
	index := 0
	nameRes,versionRes := "", ""
	for i := 0; i < len; i++ {
		_, err := strconv.Atoi(jarVersionSplit[i][0:1])
		if err == nil {
			// 是数字,几下数组下标
			index = i
			break
		}
	}
	if index == 0 {
		// 直接返回jar
		return jar, ""
	} else {
		// 把index之前的分隔数组再拼起来
		for i := 0; i < index; i++ {
			nameRes += jarVersionSplit[i]
		}

		// 把index及之后的数组拼接起来作为版本
		for i := index; i < len; i++ {
			versionRes += jarVersionSplit[i]
		}
	}
	return nameRes, versionRes
}