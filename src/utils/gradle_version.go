package utils

import (
	"strconv"
	"strings"
)

const GRADLE_PREFIX  = "gradle-"

func GradleVersion(gradleDirName string) string {
	splits := strings.Split(gradleDirName, GRADLE_PREFIX)
	version := splits[1]
	_, err := strconv.ParseFloat(version, 64)
	if err == nil {
		return version
	} else {
		return ""
	}
}