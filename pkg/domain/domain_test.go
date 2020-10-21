package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldsMapStringOK(t *testing.T) {
	timeT, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:10:10")
	ad := Ad{
		ListID:        1,
		CategoryID:    2020,
		CommuneID:     200,
		RegionID:      0,
		UserID:        111,
		Type:          "sell",
		Phone:         "999999999",
		Region:        "metropolitana",
		Commune:       "santiago",
		Category:      "vehiculos",
		Price:         100000,
		OldPrice:      200000,
		Currency:      "peso",
		ListTime:      timeT,
		PublisherType: Pro,
	}
	expected := map[string]string{
		"listid":        "1",
		"categoryid":    "2020",
		"communeid":     "200",
		"regionid":      "0",
		"userid":        "111",
		"type":          "sell",
		"phone":         "999999999",
		"region":        "metropolitana",
		"commune":       "santiago",
		"category":      "vehiculos",
		"subcategory":   "",
		"name":          "",
		"subject":       "",
		"body":          "",
		"price":         "100000",
		"oldprice":      "200000",
		"currency":      "peso",
		"listtime":      "2020-01-01 10:10:10",
		"url":           "",
		"publishertype": "pro",
	}
	result := ad.GetFieldsMapString()
	assert.Equal(t, expected, result)
}

func TestGetFieldsMapStringWithAdParamsOK(t *testing.T) {
	timeT, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:10:10")
	ad := Ad{
		ListID:        1,
		CategoryID:    2020,
		CommuneID:     200,
		RegionID:      0,
		UserID:        111,
		Type:          "sell",
		Phone:         "999999999",
		Region:        "metropolitana",
		Commune:       "santiago",
		Category:      "vehiculos",
		Price:         100000,
		OldPrice:      200000,
		Currency:      "peso",
		ListTime:      timeT,
		PublisherType: Pro,
		AdParams: map[string]string{
			"Test": "test",
			"Type": "Duplicated"},
	}
	expected := map[string]string{
		"listid":        "1",
		"categoryid":    "2020",
		"communeid":     "200",
		"regionid":      "0",
		"userid":        "111",
		"type":          "sell",
		"phone":         "999999999",
		"region":        "metropolitana",
		"commune":       "santiago",
		"category":      "vehiculos",
		"subcategory":   "",
		"name":          "",
		"subject":       "",
		"body":          "",
		"price":         "100000",
		"oldprice":      "200000",
		"currency":      "peso",
		"listtime":      "2020-01-01 10:10:10",
		"url":           "",
		"publishertype": "pro",
		"test":          "test",
	}
	result := ad.GetFieldsMapString()
	assert.Equal(t, expected, result)
}
