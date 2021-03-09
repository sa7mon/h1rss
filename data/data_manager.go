package data

import (
	"github.com/gorilla/feeds"
	"sync"
)

type DataManager struct {
	CurrentFeed *feeds.Feed
}

var once sync.Once
var instance DataManager

func GetManager() *DataManager {
	once.Do(func() { 			// atomic, do only once
		instance = DataManager{}
	})

	return &instance
}