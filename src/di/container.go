package di

import (
	"sync"
	"project/myGradle/src/jarhandler"
)

var pool = sync.Pool{
	New: func() interface{} {
		return &jarhandler.JarHandler{}
	},
}


func NewJarHandler() *jarhandler.JarHandler {
	jarH := pool.Get().(*jarhandler.JarHandler)
	if jarH == nil {
		jarH = new(jarhandler.JarHandler)
		pool.Put(jarH)
		return jarH
	}

	return jarH
}