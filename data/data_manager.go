package data

import (
	"github.com/gorilla/feeds"
	"github.com/sa7mon/h1rss/structs"
	"sync"
)

type DataManager struct {
	CurrentFeed *feeds.Feed
	ScrapedItems []structs.HacktivityItem
}

var once sync.Once
var instance DataManager

func GetManager() *DataManager {
	once.Do(func() { 			// atomic, do only once
		instance = DataManager{}
	})

	return &instance
}