package unittest

import (
	"project/myGradle/src/utils"
	"testing"
)

func Test_gradle_version(t *testing.T) {
	println(utils.GradleVersion("gradle-4.4"))
}