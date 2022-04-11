package usecases

import (
	"time"

	"github.mpi-internal.com/Yapo/ads-recommender/pkg/domain"
)

// AdMedia holds ad images data
type AdMedia struct {
	// ID image unique ID
	ID int `json:"ID"`
	// SeqNo is the image sequence number to display in inblocket platform
	SeqNo int `json:"SeqNo"`
}

// Ad holds ad response from external source
type Ad struct {
	AdID          int64                `json:"adId"`
	ListID        int64                `json:"listId"`
	UserID        int64                `json:"userId"`
	Type          string               `json:"type"`
	Phone         string               `json:"phone"`
	Location      Location             `json:"location"`
	Category      Category             `json:"category"`
	Name          string               `json:"name"`
	URL           string               `json:"url"`
	Subject       string               `json:"subject"`
	Body          string               `json:"body"`
	Price         float64              `json:"price"`
	OldPrice      float64              `json:"oldPrice"`
	ListTime      time.Time            `json:"listTime"`
	Media         []AdMedia            `json:"media"`
	PublisherType domain.PublisherType `json:"publisherType"`
	Params        map[string]Param     `json:"params"`
}

// Param represents aditional parameters on ads
type Param struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Translate interface{} `json:"translate"`
}

// Category represents a Yapo category details
type Category struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	ParentID   int64  `json:"parentId"`
	ParentName string `json:"parentName"`
}

// Location represents a location object on Ad
type Location struct {
	RegionID    int64  `json:"regionId"`
	RegionName  string `json:"regionName"`
	ComunneID   int64  `json:"communeId"`
	CommuneName string `json:"communeName"`
}

// UFApiResponse represents the indicators api response
type UFApiResponse struct {
	Version     string `json:"version"`
	Author      string `json:"autor"` // nolint: misspell
	Code        string `json:"codigo"`
	Name        string `json:"nombre"`
	MeasureUnit string `json:"unidad_medida"`
	Sets        []Set  `json:"serie"`
}

// Set represents serie in the api response
type Set struct {
	Date  string  `json:"fecha"`
	Value float64 `json:"valor"`
}

// SuggestionParameters contains all values to
// determinate which Ads should be retrieved as suggestions
type SuggestionParameters struct {
	Fields      []string
	Musts       map[string]string
	Shoulds     map[string]string
	MustsNot    map[string]string
	Filters     map[string]string
	DecayConf   map[string]string
	PriceConf   map[string]string
	QueryConf   map[string]string
	QueryString []map[string]string
}
