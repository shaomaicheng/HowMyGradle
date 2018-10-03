package model

import "container/list"

type JarCache struct {
	Name     string
	Versions list.List
}

