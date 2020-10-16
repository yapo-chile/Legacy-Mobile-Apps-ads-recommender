package usecases

import (
	"time"

	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
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
	AdID          int64                `json:"AdID"`
	ListID        int64                `json:"ListID"`
	CategoryID    int64                `json:"CategoryID"`
	CommuneID     int64                `json:"CommuneID"`
	RegionID      int64                `json:"RegionID"`
	UserID        int64                `json:"UserID"`
	Type          string               `json:"Type"`
	Phone         string               `json:"Phone"`
	Region        string               `json:"Region"`
	Commune       string               `json:"Commune"`
	Category      string               `json:"Category"`
	SubCategory   string               `json:"SubCategory"`
	Name          string               `json:"Name"`
	URL           string               `json:"Url"`
	Subject       string               `json:"Subject"`
	Body          string               `json:"Body"`
	Price         float64              `json:"Price"`
	OldPrice      float64              `json:"OldPrice"`
	ListTime      time.Time            `json:"ListTime"`
	Media         []AdMedia            `json:"Media"`
	PublisherType domain.PublisherType `json:"PublisherType"`
	Params        map[string]string    `json:"Params"`
}
