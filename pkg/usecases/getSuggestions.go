package usecases

import (
	"strconv"

	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
)

// GetSuggestions contains the repositories needed for execute a query to elastic search
type GetSuggestions struct {
	SuggestionsRepo   AdsRepository
	MinDisplayedAds   int
	RequestedAdsQty   int
	MaxDisplayedAds   int
	SuggestionsParams []string
	Logger            GetSuggestionsLogger
}

// GetSuggestionsLogger defines the logger methods that will be used for this use case
type GetSuggestionsLogger interface {
	LimitExceeded(size, maxDisplayedAds, defaultAdsQty int)
	MinimumQtyNotEnough(size, minDisplayedAds, defaultAdsQty int)
	ErrorGettingAd(listID string, err error)
	ErrorGettingAds(musts, shoulds, mustsNot map[string]string, err error)
	NotEnoughAds(listID string, lenAds int)
}

// GetProSuggestions search ad details using listId and returns a slice with ad objects
// When ad with listID is found, use the parameters on this object to search ad suggestions.
// When suggestions retrieved on repo are less than MinDisplayedAds value, itreturns empty slice.
// If something goes wrong returns empty slice and error.
func (interactor *GetSuggestions) GetProSuggestions(
	listID string, size, from int,
) (ads []domain.Ad, err error) {
	if size > interactor.MaxDisplayedAds {

		interactor.Logger.LimitExceeded(size, interactor.MaxDisplayedAds, interactor.RequestedAdsQty)
		size = interactor.RequestedAdsQty
	}
	if size < interactor.MinDisplayedAds {
		interactor.Logger.MinimumQtyNotEnough(size, interactor.MinDisplayedAds, interactor.RequestedAdsQty)
		size = interactor.RequestedAdsQty
	}

	ad, err := interactor.SuggestionsRepo.GetAd(listID)
	if err != nil {
		interactor.Logger.ErrorGettingAd(listID, err)
		return
	}
	mustParameters := getMustsParams(ad)
	shouldParameters := getShouldsParams(ad, interactor.SuggestionsParams)
	mustNotParameters := getMustNotParams(ad)
	ads, err = interactor.SuggestionsRepo.GetAds(
		mustParameters, shouldParameters, mustNotParameters, map[string]string{},
		size, from,
	)

	if err != nil {
		interactor.Logger.ErrorGettingAds(mustParameters, shouldParameters, mustNotParameters, err)
		return
	}
	// log info if there is not enough ads to return
	if len(ads) < interactor.MinDisplayedAds {
		interactor.Logger.NotEnoughAds(listID, len(ads))
		return []domain.Ad{}, nil
	}
	return ads, nil
}

func getMustsParams(ad domain.Ad) (out map[string]string) {
	out = make(map[string]string)
	out["Category"] = ad.Category
	out["SubCategory"] = ad.SubCategory
	out["PublisherType"] = string(domain.Pro)
	return
}

func getShouldsParams(ad domain.Ad, suggestionsParams []string) (out map[string]string) {
	out = make(map[string]string)
	for _, val := range suggestionsParams {
		if ad.AdParams[val] != "" {
			out["Params."+val] = ad.AdParams[val]
		}
	}
	return
}

func getMustNotParams(ad domain.Ad) (out map[string]string) {
	out = make(map[string]string)
	out["ListID"] = strconv.FormatInt(ad.ListID, 10)
	return
}
