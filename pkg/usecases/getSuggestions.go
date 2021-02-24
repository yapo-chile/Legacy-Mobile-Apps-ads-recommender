package usecases

import (
	"fmt"
	"strconv"
	"strings"

	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
)

const (
	contactField = "phonelink"

	ErrGetUF           = "ERR_GET_UF_VALUE"
	ErrInvalidCarousel = "invalid carousel: '%s'"
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
	ErrorGettingUF(err error)
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
		err = fmt.Errorf(ErrInvalidCarousel, carouselType)
		return
	}

	priceParameters, err := interactor.getPriceRange(ad, interactor.SuggestionsParams[carouselType]["priceRange"])
	if err != nil {
		interactor.Logger.ErrorGettingUF(err)
		return
	}
	decayParameters := getDecayFunctionParams(interactor.SuggestionsParams, carouselType)
	queryStringParameters := getQueryStringParams(interactor.SuggestionsParams[carouselType]["queryString"])
	mustParameters := getMustsParams(ad, interactor.SuggestionsParams[carouselType]["must"])
	shouldParameters := getShouldsParams(ad, interactor.SuggestionsParams[carouselType]["should"])
	mustNotParameters := getMustNotParams(ad, interactor.SuggestionsParams[carouselType]["mustNot"])
	filtersParameters := getFilterParams(ad, interactor.SuggestionsParams[carouselType]["filter"])

	ads, err = interactor.SuggestionsRepo.GetAds(
		mustParameters, shouldParameters, mustNotParameters, filtersParameters,
		priceParameters,
		decayParameters,
		queryStringParameters,
		size, from,
	)

	if err != nil {
		interactor.Logger.ErrorGettingAds(mustParameters, shouldParameters, mustNotParameters, err)
		return
	}

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

// getPriceRange returns a map with price range values
func (interactor *GetSuggestions) getPriceRange(
	ad domain.Ad,
	priceRangeSlice []interface{},
) (out map[string]string, err error) {
	out = make(map[string]string)
	if len(priceRangeSlice) <= 0 {
		return
	}

	uf, err := interactor.IndicatorsRepository.GetUF()
	if err != nil {
		return
	}

	priceRange := priceRangeSlice[0].(map[string]interface{})

	out["uf"] = fmt.Sprintf("%v", uf)
	if _, ok := priceRange["type"]; !ok {
		out["type"] = "must"
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

// getQueryStringParams returns a map query string values
func getQueryStringParams(queryStringSlice []interface{}) (out []map[string]string) {
	out = make([]map[string]string, 0)
	if len(queryStringSlice) <= 0 {
		return
	}

	for _, value := range queryStringSlice {
		valueStr := value.(map[string]interface{})
		outTmp := make(map[string]string)
		outTmp["query"] = fmt.Sprintf("%v", valueStr["query"])
		outTmp["default_field"] = fmt.Sprintf("%v", valueStr["defaultField"])

		out = append(out, outTmp)
	}
	return
}

// getMustsParams returns a map with mandatory values
func getMustsParams(ad domain.Ad, suggestionsParams []interface{}) (out map[string]string) {
	adMap := ad.GetFieldsMapString()
	out = getParams(adMap, suggestionsParams)
	return
}

// getShouldsParams returns a map with optional values
func getShouldsParams(ad domain.Ad, suggestionsParams []interface{}) (out map[string]string) {
	adMap := ad.GetFieldsMapString()
	return getParams(adMap, suggestionsParams)
}

// getMustNotParams returns a map with excluded values
func getMustNotParams(ad domain.Ad, suggestionsParams []interface{}) (out map[string]string) {
	adMap := ad.GetFieldsMapString()
	return getParams(adMap, suggestionsParams)
}

// getFilterParams returns a map with filter values which dont add score
// to the resulting documents
func getFilterParams(ad domain.Ad, suggestionsParams []interface{}) (out map[string]string) {
	adMap := ad.GetFieldsMapString()
	return getParams(adMap, suggestionsParams)
}

// getDecayFunctionParams returns a map with decay function values
func getDecayFunctionParams(decayFuncConf map[string]map[string][]interface{}, carouselType string) (out map[string]string) {
	out = make(map[string]string)
	if len(decayFuncConf[carouselType]["decayFunc"]) <= 0 {
		carouselType = "default"
	}

	decayFunc := decayFuncConf[carouselType]["decayFunc"][0].(map[string]interface{})

	out["name"] = decayFunc["name"].(string)
	out["field"] = decayFunc["field"].(string)
	out["origin"] = decayFunc["origin"].(string)
	out["offset"] = decayFunc["offset"].(string)
	out["scale"] = decayFunc["scale"].(string)

	return
}

// calculateMinMaxPriceRange calculates the minimum and maximum
// price for a range query, where a value is substracted and added to the price
// of the ad being requested
func calculateMinMaxPriceRange(
	adPrice, uf float64,
	adCurrency string,
	minusValue, plusValue int,
) (string, string) {
	// if ad currency is 'peso', divide price by UF to convert to equivalent
	// UF value
	if adCurrency == "peso" {
		adPrice /= uf
	} else {
		adPrice /= 100
	}

	minPrice := adPrice - float64(minusValue)
	maxPrice := adPrice + float64(plusValue)

	return fmt.Sprintf("%v", minPrice), fmt.Sprintf("%v", maxPrice)
}

// getParams function that reads the parameters to be used in the queries
// must, mustNot, should and filter, where they can come in the string form
// Params.{param} or simply as {param}
func getParams(adMap map[string]string, suggestionsParams []interface{}) (out map[string]string) {
	out = make(map[string]string)
	for _, param := range suggestionsParams {
		paramKey := param.(string)
		var paramValue string

		if strings.HasPrefix(paramKey, "Params.") {
			paramSlice := strings.Split(paramKey, ".")
			paramValue = strings.ToLower(paramSlice[len(paramSlice)-1])
		} else {
			paramValue = strings.ToLower(paramKey)
		}

		if adMap[paramValue] != "" {
			out[paramKey] = adMap[paramValue]
		}
	}
	return
}
