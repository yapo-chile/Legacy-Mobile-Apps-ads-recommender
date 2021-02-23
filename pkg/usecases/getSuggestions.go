package usecases

import (
	"fmt"
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
	InvalidCarousel(carousel string)
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

	if _, ok := interactor.SuggestionsParams[carouselType]; !ok {
		interactor.Logger.InvalidCarousel(carouselType)
		err = fmt.Errorf("invalid carousel: '%s'", carouselType)
		return
	}

	priceParameters := interactor.getPriceRange(ad, interactor.SuggestionsParams[carouselType]["priceRange"])
	mustParameters := getMustsParams(ad, interactor.SuggestionsParams[carouselType]["must"])
	shouldParameters := getShouldsParams(ad, interactor.SuggestionsParams[carouselType]["should"])
	mustNotParameters := getMustNotParams(ad, interactor.SuggestionsParams[carouselType]["mustNot"])
	filtersParameters := getFilterParams(ad, interactor.SuggestionsParams[carouselType]["filter"])

	ads, err = interactor.SuggestionsRepo.GetAds(
		mustParameters, shouldParameters, mustNotParameters, filtersParameters,
		priceParameters,
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
	suggestionsParams []interface{},
) (out map[string]string) {
	out = make(map[string]string)
	out["PublisherType"] = string(domain.Pro)
	adMap := ad.GetFieldsMapString()

	for _, val := range suggestionsParams {
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
	if len(priceRangeSlice) <= 0 {
		return
	}

	uf, _ := interactor.IndicatorsRepository.GetUF()
	priceRange := priceRangeSlice[0].(map[string]interface{})

	out["uf"] = fmt.Sprintf("%v", uf)
	if _, ok := priceRange["type"]; !ok {
		priceRange["type"] = "must"
	} else {
		out["type"] = priceRange["type"].(string)
	}

	if _, ok := priceRange["calculate"]; ok {
		adMap := ad.GetFieldsMapString()

		adCurrency := adMap["currency"]
		adPrice, _ := strconv.ParseFloat(adMap["price"], 64)
		minusPrice, _ := strconv.Atoi(priceRange["gte"].(string))
		plusPrice, _ := strconv.Atoi(priceRange["lte"].(string))

		out["gte"], out["lte"] =
			calculateMinMaxPriceRange(adPrice, uf, adCurrency, minusPrice, plusPrice)
	} else {
		out["gte"], out["lte"] = priceRange["gte"].(string), priceRange["lte"].(string)
	}
	return
}

// getShouldsParams returns a map with optional values
func getShouldsParams(
	ad domain.Ad,
	suggestionsParams []interface{},
) (out map[string]string) {
	out = make(map[string]string)
	adMap := ad.GetFieldsMapString()

	//should params can come with syntax 'Params.{param} or simply {param}'
	//se we need to split the string with '.' in the first syntax is used to
	//get the value
	for _, shoulParam := range suggestionsParams {
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
func getMustNotParams(ad domain.Ad, suggestionsParams []interface{}) (out map[string]string) {
	out = make(map[string]string)
	adMap := ad.GetFieldsMapString()

	for _, val := range suggestionsParams {
		v := val.(string)
		if adMap[strings.ToLower(v)] != "" {
			out[v] = adMap[strings.ToLower(v)]
		}
	}
	return
}

// getFilterParams returns a map with filter values
func getFilterParams(ad domain.Ad, suggestionsParams []interface{}) (out map[string]string) {
	out = make(map[string]string)
	adMap := ad.GetFieldsMapString()

	for _, val := range suggestionsParams {
		v := val.(string)
		if adMap[strings.ToLower(v)] != "" {
			out[v] = adMap[strings.ToLower(v)]
		}
	}
	return
}

// getDecayFunctionParams returns a map with decay function values
func getDecayFunctionParams(suggestionsParams []interface{}) (out map[string]string) {
	out = make(map[string]string)

	defaultName := "gauss"
	defaultField := "ListTime"
	defaultOrigin := "now/1d"
	defaultOffset := "1d"
	defaultScale := "7d"

	if len(suggestionsParams) < 0 {
		out["name"] = defaultName
		out["field"] = defaultField
		out["origin"] = defaultOrigin
		out["offset"] = defaultOffset
		out["scale"] = defaultScale
	}

	decayFuncParams := suggestionsParams[0].(map[string]interface{})

	if val, ok := decayFuncParams["name"]; ok {
		out["name"] = val.(string)
	} else {
		out["name"] = defaultName
	}

	if val, ok := decayFuncParams["field"]; ok {
		out["field"] = val.(string)
	} else {
		out["field"] = defaultField
	}

	if val, ok := decayFuncParams["origin"]; ok {
		out["origin"] = val.(string)
	} else {
		out["origin"] = defaultOrigin
	}

	if val, ok := decayFuncParams["offset"]; ok {
		out["offset"] = val.(string)
	} else {
		out["offset"] = defaultOffset
	}

	if val, ok := decayFuncParams["scale"]; ok {
		out["scale"] = val.(string)
	} else {
		out["scale"] = defaultScale
	}

	return
}

func calculateMinMaxPriceRange(
	adPrice, uf float64,
	adCurrency string,
	minusValue, plusValue int,
) (string, string) {
	// if ad currency is 'peso', divide price by UF
	if adCurrency == "peso" {
		adPrice /= uf
	} else {
		adPrice /= 100
	}

	minPrice := adPrice - float64(minusValue)
	maxPrice := adPrice + float64(plusValue)

	return fmt.Sprintf("%v", minPrice), fmt.Sprintf("%v", maxPrice)
}
