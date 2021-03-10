package structs

import (
	"github.com/gorilla/feeds"
	"time"
)

type HacktivityItem struct {
	Title string
	State string
	Link string
	Bounty string
	HasBounty bool
	Severity string
	ReportedBy string
	ReportedTo string
	LastUpdate time.Time
	RSSItem *feeds.Item
}

type H1GraphResponse struct {
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
