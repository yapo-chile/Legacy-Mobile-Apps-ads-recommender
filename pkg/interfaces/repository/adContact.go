package repository

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.mpi-internal.com/Yapo/ads-recommender/pkg/domain"
	"github.mpi-internal.com/Yapo/ads-recommender/pkg/usecases"
)

const (
	errorRetrievingAds = "there was an error retrieving ads info: %+v"
	errorNoAdsFound    = "the specified ads %+v don't return results from ad-contact"
	errorUnMarshallAds = "there was an error retrieving ads phone list from ad-contact: %+v. \nBody request: %+v"
)

// AdsContactPhonesInput object to set input to get ad contact info
type adsContactPhonesInput struct {
	ListIDs []string `json:"list_ids"`
}

// AdContactRepository implements the handler to get ad's contact info
type AdContactRepository struct {
	handler HTTPHandler
	path    string
}

// NewAdContactRepository returns a fresh instance of AdContactRepo
func NewAdContactRepository(handler HTTPHandler, path string) usecases.AdContactRepo {
	return &AdContactRepository{
		handler: handler,
		path:    path,
	}
}

// GetAdsPhone gets ads contact info from ad contact ms
func (repo *AdContactRepository) GetAdsPhone(
	suggestions []domain.Ad,
) (adsResult map[string]string, err error) {
	var listIds []string
	for _, ad := range suggestions {
		listIds = append(listIds, strconv.FormatInt(ad.ListID, 10))
	}
	request := repo.handler.NewRequest().
		SetMethod("GET").
		SetPath(repo.path).
		SetBody(adsContactPhonesInput{ListIDs: listIds})
	adsJSON, err := repo.handler.Send(request)
	if err == nil && adsJSON != nil {
		ads := fmt.Sprintf("%s", adsJSON)
		err = json.Unmarshal([]byte(ads), &adsResult)
		if err != nil {
			return adsResult, fmt.Errorf(errorUnMarshallAds, err, ads)
		}
		if len(adsResult) == 0 {
			return adsResult, fmt.Errorf(errorNoAdsFound, ads)
		}
		return adsResult, nil
	}
	return adsResult, fmt.Errorf(errorRetrievingAds, err)
}
