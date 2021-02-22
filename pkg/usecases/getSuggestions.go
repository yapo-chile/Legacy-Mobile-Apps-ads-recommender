package usecases

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
)

const (
	contactField = "phonelink"

	ErrGetUF = "ERR_GET_UF_VALUE"
)

// GetSuggestions contains the repositories needed to retrieve ads suggestions
type GetSuggestions struct {
	SuggestionsRepo      AdsRepository
	AdContact            AdContactRepo
	MinDisplayedAds      int
	RequestedAdsQty      int
	MaxDisplayedAds      int
	SuggestionsParams    map[string]map[string][]interface{}
	Logger               GetSuggestionsLogger
	IndicatorsRepository IndicatorsRepository
}

// GetSuggestionsLogger defines the logger methods that will be used for this usecase
type GetSuggestionsLogger interface {
	LimitExceeded(size, maxDisplayedAds, defaultAdsQty int)
	MinimumQtyNotEnough(size, minDisplayedAds, defaultAdsQty int)
	ErrorGettingAd(listID string, err error)
	ErrorGettingAds(musts, shoulds, mustsNot map[string]string, err error)
	NotEnoughAds(listID string, lenAds int)
	ErrorGettingAdsContact(listID string, err error)
}

// GetProSuggestions search ad details using listId and returns a slice with ad objects
// When ad with listID is found, use the parameters on this object to search ad suggestions.
// When suggestions retrieved on repo are less than MinDisplayedAds value, it returns empty slice.
// If something goes wrong returns empty slice and error.
func (interactor *GetSuggestions) GetProSuggestions(
	listID string, optionalParams []string, size, from int, carouselType string,
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

	if carouselType == "" {
		carouselType = "default"
	}

	if _, ok := interactor.SuggestionsParams[carouselType]; !ok {
		return
	}

	rangeParameters := interactor.getPriceRange(ad, interactor.SuggestionsParams[carouselType]["priceRange"])
	log.Printf("rangeParameters %v", rangeParameters)
	mustParameters := getMustsParams(ad, interactor.SuggestionsParams, carouselType)
	shouldParameters := getShouldsParams(ad, interactor.SuggestionsParams, carouselType)
	mustNotParameters := getMustNotParams(ad)
	ads, err = interactor.SuggestionsRepo.GetAds(
		mustParameters, shouldParameters, mustNotParameters, map[string]string{},
		rangeParameters,
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
	ads, err = interactor.getAdsContact(ads, optionalParams)
	if err != nil {
		interactor.Logger.ErrorGettingAdsContact(listID, err)
	}
	return ads, nil
}

// getAdsContact if phonelink is required connect to adContact repo
// and gets ads contact data.
func (interactor *GetSuggestions) getAdsContact(
	suggestions []domain.Ad,
	optionalParams []string,
) (ads []domain.Ad, err error) {
	phones := make(map[string]string)
	for _, param := range optionalParams {
		if strings.EqualFold(param, contactField) {
			phones, err = interactor.AdContact.GetAdsPhone(suggestions)
			break
		}
	}
	if len(phones) > 0 {
		for _, ad := range suggestions {
			if val, ok := phones[strconv.FormatInt(ad.ListID, 10)]; ok {
				if ad.AdParams == nil {
					ad.AdParams = make(map[string]string)
				}
				ad.AdParams[contactField] = val
			}
		}
	}
	return suggestions, err
}

// getMustsParams returns a map with mandatory values
func getMustsParams(
	ad domain.Ad,
	suggestionsParams map[string]map[string][]interface{},
	carouselType string,
) (out map[string]string) {
	out = make(map[string]string)
	out["PublisherType"] = string(domain.Pro)
	adMap := ad.GetFieldsMapString()

	for _, val := range suggestionsParams[carouselType]["must"] {
		v := val.(string)
		if adMap[strings.ToLower(v)] != "" {
			out[v] = adMap[strings.ToLower(v)]
		}
	}
	return
}

// getPriceRange returns a map with range values
func (interactor *GetSuggestions) getPriceRange(
	ad domain.Ad,
	priceRangeSlice []interface{},
) (out map[string]string) {
	out = make(map[string]string)
	uf, _ := interactor.IndicatorsRepository.GetUF()
	priceRange := priceRangeSlice[0].(map[string]interface{})
	for key, val := range priceRange {
		out[key] = val.(string)
		if key == "calculate" {
			adMap := ad.GetFieldsMapString()

			// if ad currency is 'peso', divide price by UF
			var adPrice float64
			if adMap["currency"] == "peso" {
				adPrice, _ = strconv.ParseFloat(adMap["price"], 64)
				adPrice /= uf

			}

			minusPrice, _ := strconv.Atoi(priceRange["gte"].(string))
			plusPrice, _ := strconv.Atoi(priceRange["lte"].(string))
			minPrice := adPrice - float64(minusPrice)
			maxPrice := adPrice + float64(plusPrice)

			out["gte"] = fmt.Sprintf("%v", minPrice)
			out["lte"] = fmt.Sprintf("%v", maxPrice)
		}
	}
	out["uf"] = fmt.Sprintf("%v", uf)
	return
}

// getShouldsParams returns a map with optional values
func getShouldsParams(
	ad domain.Ad,
	suggestionsParams map[string]map[string][]interface{},
	carouselType string,
) (out map[string]string) {
	out = make(map[string]string)
	adMap := ad.GetFieldsMapString()

	//should params can come with syntax 'Params.{param} or simply {param}'
	//se we need to split the string with '.' in the first syntax is used to
	//get the value
	for _, shoulParam := range suggestionsParams[carouselType]["should"] {
		shoulParamKey := shoulParam.(string)
		shoulParamSlice := strings.Split(shoulParamKey, ".")
		shouldParamValue := shoulParamSlice[len(shoulParamSlice)-1]

		if adMap[strings.ToLower(shouldParamValue)] != "" {
			out[shoulParamKey] = adMap[strings.ToLower(shouldParamValue)]
		}
	}
	return
}

// getMustNotParams returns a map with excluded values
func getMustNotParams(ad domain.Ad) (out map[string]string) {
	out = make(map[string]string)
	out["ListID"] = strconv.FormatInt(ad.ListID, 10)
	return
}

func caclulateMinMaxPriceRange(adPrice float64, minusValue, plusValue int) (float64, error) {
	return 0, nil
}
