package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Ad struct for single ad representation
type Ad struct {
	AdID             int64
	ListID           int64
	CategoryID       int64
	CategoryParentID int64
	CommuneID        int64
	RegionID         int64
	UserID           int64
	Type             string
	Phone            string
	Region           string
	Commune          string
	Category         string
	CategoryParent   string
	Name             string
	Subject          string
	Body             string
	Price            float64
	OldPrice         float64
	Currency         string
	ListTime         time.Time
	URL              string
	Image            Image
	PublisherType    PublisherType
	AdParams         map[string]string
}

// GetFieldsMapString returns a map with all fields and values
func (ad *Ad) GetFieldsMapString() (output map[string]string) {
	output = map[string]string{
		"listid":           strconv.FormatInt(ad.ListID, 10),
		"categoryid":       strconv.FormatInt(ad.CategoryID, 10),
		"communeid":        strconv.FormatInt(ad.CommuneID, 10),
		"regionid":         strconv.FormatInt(ad.RegionID, 10),
		"userid":           strconv.FormatInt(ad.UserID, 10),
		"type":             ad.Type,
		"phone":            ad.Phone,
		"region":           ad.Region,
		"commune":          ad.Commune,
		"category":         ad.Category,
		"categoryparent":   ad.CategoryParent,
		"categoryparentid": strconv.FormatInt(ad.CategoryParentID, 10),
		"name":             ad.Name,
		"subject":          ad.Subject,
		"body":             ad.Body,
		"price":            fmt.Sprintf("%g", ad.Price),
		"oldprice":         fmt.Sprintf("%g", ad.OldPrice),
		"currency":         ad.Currency,
		"listtime":         ad.ListTime.Format("2006-01-02 15:04:05"),
		"url":              ad.URL,
		"publishertype":    string(ad.PublisherType),
	}
	for key, val := range ad.AdParams {
		key = strings.ToLower(key)
		if output[key] == "" {
			output[key] = val
		}
	}
	return
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
