package loggers

import (
	"fmt"
	"testing"
)

func TestGetSuggestionsLogger(t *testing.T) {
	m := &loggerMock{t: t}
	mMap := map[string]string{}
	l := MakeGetSuggestionsLogger(m)
	l.LimitExceeded(0, 0, 0)
	l.MinimumQtyNotEnough(0, 0, 0)
	l.ErrorGettingAd("", fmt.Errorf(""))
	l.ErrorGettingAds(mMap, mMap, mMap, fmt.Errorf(""))
	l.NotEnoughAds("", 0)
	l.ErrorGettingAdsContact("", fmt.Errorf(""))
	m.AssertExpectations(t)
}
