package unittest

import (
	"regexp"
	"testing"
)

func Test_regex(t *testing.T) {
	regex := "gradle-[1-9]+.[1-9]+"
	match1, _ := regexp.MatchString(regex, "gradle-4.4")
	println(match1)
	match2, _ := regexp.MatchString(regex, "gradle")
	println(match2)
	match3, _ := regexp.MatchString(regex, "gradle-six6.6")
	println(match3)
}

func Test_jar_regex(t *testing.T) {
	jarRegex := "(.*).jar"
	match1, _ := regexp.MatchString(jarRegex, "baseLibrary-3.0.1.jar")
	println(match1)
}
