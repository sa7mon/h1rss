package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/sa7mon/h1rss/data"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var scrapeInterval int
	var bindAddr string
	flag.IntVar(&scrapeInterval, "interval", 120, "Minutes to wait between scrapes")
	flag.StringVar(&bindAddr, "bind", ":8000", "Address and port to bind to")
	flag.Parse()

	if scrapeInterval < 1 {
		fmt.Println("Scraping interval must be at least 1")
		os.Exit(1)
	}

	if !strings.ContainsRune(bindAddr, ':') || len(bindAddr) < 2 {
		fmt.Println("flag 'bind' must be in format 'address:port'. If address is omitted, server will listed on all interfaces.")
		os.Exit(1)
	}

	log.Printf("[main] Scraping hacktivity every %v minutes", scrapeInterval)
	manager := data.GetManager()

	s := NewScraper()
	rssItems, err := s.Scrape()
	if err != nil {
		panic(err)
	}
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "HackerOne Unofficial Hacktivity RSS Feed",
		Link:        &feeds.Link{Href: ""},
		Description: "",
		Author:      &feeds.Author{Name: "", Email: ""},
		Created:     now,
		Items: rssItems,
	}

	manager.CurrentFeed = feed

	r := mux.NewRouter()
	r.HandleFunc("/rss", RSSHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         bindAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Spin off scraper to its own thread
	go s.ScrapeLoop(2)

	log.Printf("[server] Serving on %v", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func RSSHandler(w http.ResponseWriter, r *http.Request) {
	manager := data.GetManager()
	rss, err := manager.CurrentFeed.ToRss()
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/rss+xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rss))
}

type scraper struct {
	feedItems []feeds.Item
}

func NewScraper() scraper {
	return scraper{}
}

func (sc scraper) ScrapeLoop(interval int) {
	keepScraping := true
	sleepInterval := time.Duration(interval) * time.Minute * 60
	manager := data.GetManager()
	var scrapeError error

	for keepScraping {
		time.Sleep(sleepInterval)
		log.Printf("[scraper] Starting scrape")
		items, err := sc.Scrape()
		if err != nil {
			scrapeError = err
			keepScraping = false
			continue
		}
		manager.CurrentFeed.Items = items
	}
	log.Printf("[scraper] Thread dying due to error: %v", scrapeError)
}

/*
	Adapted from: https://gist.github.com/tetrillard/4e1ed77cebb5fab42989da3bf944fd4e
 */
func (sc scraper) Scrape() ([]*feeds.Item, error) {
	data := `{"operationName":"HacktivityPageQuery","variables":{"querystring":"","where":{"report":{"disclosed_at":{"_is_null":false}}},"orderBy":null,"secureOrderBy":{"latest_disclosable_activity_at":{"_direction":"DESC"}},"count":10},"query":"query HacktivityPageQuery($querystring: String, $orderBy: HacktivityItemOrderInput, $secureOrderBy: FiltersHacktivityItemFilterOrder, $where: FiltersHacktivityItemFilterInput, $count: Int, $cursor: String) {\n  hacktivity_items(first: $count, after: $cursor, query: $querystring, order_by: $orderBy, secure_order_by: $secureOrderBy, where: $where) {\n    ...HacktivityList\n  }\n}\n\nfragment HacktivityList on HacktivityItemConnection {\n    edges {\n    node {\n      ... on HacktivityItemInterface {\n        ...HacktivityItem\n      }\n    }\n  }\n}\n\nfragment HacktivityItem on HacktivityItemUnion {\n  ... on Undisclosed {\n    id\n    ...HacktivityItemUndisclosed\n  }\n  ... on Disclosed {\n    ...HacktivityItemDisclosed\n  }\n  ... on HackerPublished {\n    ...HacktivityItemHackerPublished\n  }\n}\n\nfragment HacktivityItemUndisclosed on Undisclosed {\n  reporter {\n    username\n    ...UserLinkWithMiniProfile\n  }\n  team {\n    handle\n    name\n     url\n    ...TeamLinkWithMiniProfile\n  }\n  latest_disclosable_action\n  latest_disclosable_activity_at\n  requires_view_privilege\n  total_awarded_amount\n  currency\n}\n\nfragment TeamLinkWithMiniProfile on Team {\n  handle\n  name\n }\n\nfragment UserLinkWithMiniProfile on User {\n  username\n}\n\nfragment HacktivityItemDisclosed on Disclosed {\n  reporter {\n    username\n    ...UserLinkWithMiniProfile\n  }\n  team {\n    handle\n    name\n    url\n    ...TeamLinkWithMiniProfile\n  }\n  report {\n    title\n    substate\n    url\n  }\n  latest_disclosable_activity_at\n  total_awarded_amount\n  severity_rating\n  currency\n}\n\nfragment HacktivityItemHackerPublished on HackerPublished {\n  reporter {\n    username\n    ...UserLinkWithMiniProfile\n  }\n  team {\n    handle\n    name\n    medium_profile_picture: profile_picture(size: medium)\n    url\n    ...TeamLinkWithMiniProfile\n  }\n  report {\n    url\n    title\n    substate\n  }\n  latest_disclosable_activity_at\n  severity_rating\n}\n"}`
	data = strings.Replace(data, "\n", "\\n", -1)

	request, err := http.NewRequest("POST", "https://hackerone.com/graphql", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", response.StatusCode, response.Status))
	}
	var resp h1GraphResponse

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		return nil, err
	}

	var parsedItems []*feeds.Item
	for _, v := range resp.Data.HacktivityItems.Edges {
		var item feeds.Item
		var title string

		if v.Node.TotalAwardedAmount != 0.00 {
			title = fmt.Sprintf("[%v] [%v %v]", v.Node.Team.Name, v.Node.TotalAwardedAmount, v.Node.Currency)
		} else {
			title = fmt.Sprintf("[%v]", v.Node.Team.Name)
		}

		severity := fmt.Sprintf("%v", v.Node.SeverityRating)
		if severity != "none" && severity != "<nil>" {
			title = fmt.Sprintf("%v [%v]", title, v.Node.SeverityRating)
		}

		title = fmt.Sprintf("%v %v", title, v.Node.Report.Title)

		item = feeds.Item{
			Title: title,
			Updated: v.Node.LatestDisclosableActivityAt,
			Link: 	&feeds.Link{Href: v.Node.Report.URL},
			Description: "",
			Author: &feeds.Author{Name: "", Email: ""},
		}
		parsedItems = append(parsedItems, &item)
	}

	return parsedItems, nil
}

type h1GraphResponse struct {
	Data struct {
		HacktivityItems struct {
			Edges []struct {
				Node struct {
					Reporter struct {
						Username string `json:"username"`
					} `json:"reporter"`
					Team struct {
						Handle string `json:"handle"`
						Name   string `json:"name"`
						URL    string `json:"url"`
					} `json:"team"`
					Report struct {
						Title    string `json:"title"`
						Substate string `json:"substate"`
						URL      string `json:"url"`
					} `json:"report"`
					LatestDisclosableActivityAt time.Time      `json:"latest_disclosable_activity_at"`
					TotalAwardedAmount          float64        `json:"total_awarded_amount"`
					SeverityRating              interface{} `json:"severity_rating"`
					Currency                    string         `json:"currency"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"hacktivity_items"`
	} `json:"data"`
}