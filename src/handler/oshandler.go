package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

}