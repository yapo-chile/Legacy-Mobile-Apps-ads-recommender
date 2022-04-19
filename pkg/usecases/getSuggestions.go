package usecases

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/domain"
)

const (
	contactField = "phonelink"
	// ErrGetUF error code when get uf fails
	ErrGetUF = "ERR_GET_UF_VALUE"
	// ErrInvalidCarousel error text when an invalid carousel is requested
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

// GetSuggestions search ad details using listId and returns a slice with ad objects
// When sourceAd parameter is true, it retrieves an ad using a listID.
// It translates data from conf y/o ad fields as parameters to search a slice with ad suggestions.
// When suggestions retrieved on repo are less than MinDisplayedAds value, it returns empty slice.
// If something goes wrong returns empty slice and error.
func (interactor *GetSuggestions) GetSuggestions(
	listID string, optionalParams []string, size, from int, carouselType string,
) (ads []domain.Ad, err error) {
	ads = []domain.Ad{}
	size = interactor.getSize(size)
	if _, ok := interactor.SuggestionsParams[carouselType]; !ok {
		interactor.Logger.InvalidCarousel(carouselType)
		err = fmt.Errorf(ErrInvalidCarousel, carouselType)
		return
	}
	parameters, adID, err := interactor.getSuggestionParameters(listID, carouselType)
	if err != nil {
		return
	}

	ads, err = interactor.SuggestionsRepo.GetAds(
		adID,
		parameters,
		size,
		from,
	)
	if err != nil {
		interactor.Logger.ErrorGettingAds(
			parameters.Musts, parameters.Shoulds, parameters.MustsNot, err)
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

// getSuggestionParameters creates and retrieves a struct containing all parameters to get ad suggestions
// if something goes wrong it retrieves and empty struct and error
func (interactor *GetSuggestions) getSuggestionParameters(
	listID, carouselType string) (params SuggestionParameters, adID string, err error) {
	var ad domain.Ad
	ad, err = interactor.SuggestionsRepo.GetAd(listID)
	if err != nil {
		interactor.Logger.ErrorGettingAd(listID, err)
		return
	}
	adID = strconv.FormatInt(ad.AdID, 10)
	adMap := ad.GetFieldsMapString()
	params.PriceConf, err = interactor.getPriceRange(ad, interactor.SuggestionsParams[carouselType]["priceRange"])
	if err != nil {
		interactor.Logger.ErrorGettingUF(err)
		return
	}

	params.QueryConf = getValues(interactor.SuggestionsParams, carouselType, "queryConf")
	params.DecayConf = getValues(interactor.SuggestionsParams, carouselType, "decayFunc")
	params.QueryString = getQueryStringParams(interactor.SuggestionsParams[carouselType]["queryString"])

	params.Musts = getSliceParams(adMap, interactor.SuggestionsParams[carouselType]["must"])

	params.Shoulds = getSliceParams(adMap, interactor.SuggestionsParams[carouselType]["should"])
	params.MustsNot = getSliceParams(adMap, interactor.SuggestionsParams[carouselType]["mustNot"])
	params.Filters = getSliceParams(adMap, interactor.SuggestionsParams[carouselType]["filter"])
	params.Fields = getSliceString(interactor.SuggestionsParams[carouselType]["fields"])
	return
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
	if len(priceRangeSlice) == 0 {
		return out, err
	}

	uf, errUF := interactor.IndicatorsRepository.GetUF()
	if errUF != nil {
		interactor.Logger.ErrorGettingUF(errUF)
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

		out["gte"], out["lte"] = calculateMinMaxPriceRange(
			adPrice,
			uf,
			adCurrency,
			minusPrice,
			plusPrice,
		)
	} else {
		out["gte"], out["lte"] = priceRange["gte"].(string), priceRange["lte"].(string)
	}
	return out, err
}

// getSize retrieves default size if input size equals zero, otherwise returns size
func (interactor *GetSuggestions) getSize(size int) int {
	if size > interactor.MaxDisplayedAds {
		interactor.Logger.LimitExceeded(size, interactor.MaxDisplayedAds, interactor.RequestedAdsQty)
		size = interactor.RequestedAdsQty
	}
	if size < interactor.MinDisplayedAds {
		interactor.Logger.MinimumQtyNotEnough(size, interactor.MinDisplayedAds, interactor.RequestedAdsQty)
		size = interactor.RequestedAdsQty
	}
	return size
}

// getQueryStringParams returns a map query string values
func getQueryStringParams(queryStringSlice []interface{}) (out []map[string]string) {
	out = make([]map[string]string, 0)
	if len(queryStringSlice) == 0 {
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

// calculateMinMaxPriceRange calculates the minimum and maximum
// price for a range query, where a value is subtracted and added to the price
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
	}
	minPrice := adPrice - float64(minusValue)
	maxPrice := adPrice + float64(plusValue)

	return fmt.Sprintf("%v", minPrice), fmt.Sprintf("%v", maxPrice)
}

// getSliceParams function that reads the parameters to be used in the queries
// must, mustNot, should and filter, where they can come in the string form
// Params.{param} or simply as {param}
func getSliceParams(adMap map[string]string, suggestionsParams []interface{}) (out map[string]string) {
	out = make(map[string]string)

	for _, param := range suggestionsParams {
		paramKey := param.(string)
		var paramValue string

		if strings.HasPrefix(paramKey, "params.") || strings.HasPrefix(paramKey, "location.") {
			paramSlice := strings.Split(paramKey, ".")
			paramValue = strings.ToLower(paramSlice[1])
		} else if strings.HasPrefix(paramKey, "category.") {
			paramValue = strings.ReplaceAll(paramKey, ".", "")
		} else {
			paramValue = strings.ToLower(paramKey)
		}

		if adMap[paramValue] != "" {
			out[paramKey] = adMap[paramValue]
		}
	}
	return
}

// getSliceString transforms interface slice to string slice
func getSliceString(input []interface{}) (output []string) {
	for _, value := range input {
		if str, ok := value.(string); ok {
			output = append(output, str)
		}
	}
	return
}

// getValues returns a map with a config used in a specific carousel
func getValues(
	confValues map[string]map[string][]interface{},
	carouselType, confName string,
) (output map[string]string) {
	output = make(map[string]string)
	if len(confValues[carouselType][confName]) == 0 {
		if len(confValues["default"][confName]) == 0 {
			return
		}
		carouselType = "default"
	}
	conf := confValues[carouselType][confName][0].(map[string]interface{})
	for key, value := range conf {
		if value != nil {
			if str, ok := value.(string); ok {
				output[key] = str
			}
		}
	}
	return
}
