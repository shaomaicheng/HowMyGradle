package unittest

import (
	"testing"
	"project/myGradle/src/di"
)

func Test_get_instance(t *testing.T) {
	jarHandler1 := di.NewJarHandler()
	jarHandler2 := di.NewJarHandler()
	jarHandler3 := di.NewJarHandler()
	jarHandler4 := di.NewJarHandler()
	jarHandler5 := di.NewJarHandler()
	println(jarHandler1)
	println(jarHandler2)
	println(jarHandler3)
	println(jarHandler4)
	println(jarHandler5)
}