package domain

import (
	"time"
)

// Ad struct for single ad representation
type Ad struct {
	ListID        int64
	CategoryID    int64
	CommuneID     int64
	RegionID      int64
	UserID        int64
	Type          string
	Phone         string
	Region        string
	Commune       string
	Category      string
	SubCategory   string
	Name          string
	Subject       string
	Body          string
	Price         float64
	OldPrice      float64
	Currency      string
	ListTime      time.Time
	URL           string
	Image         Image
	PublisherType PublisherType
	AdParams      map[string]string
}

// Image struct that defines the internal structure of ad images
type Image struct {
	Full   string
	Medium string
	Small  string
}

// PublisherType describes publisher user
type PublisherType string

// Pri and Pro values
const (
	Pro PublisherType = "pro"
	Pri PublisherType = "private"
)
