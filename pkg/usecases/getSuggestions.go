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
		log.Printf("carousel \"%s\" is not a valid carousel", carouselType)
		return
	}

	rangeParameters := interactor.getRange(ad, interactor.SuggestionsParams, carouselType)
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

// getRange returns a map with range values
func (interactor *GetSuggestions) getRange(
	ad domain.Ad,
	suggestionsParams map[string]map[string][]interface{},
	carouselType string,
) (out map[string]map[string]string) {
	out = make(map[string]map[string]string)
	for _, val := range suggestionsParams[carouselType]["range"] {
		v := val.(map[string]interface{})
		for rangeKey, lim := range v {
			rng := make(map[string]string)

			rangeValues := lim.(map[string]interface{})
			for lk, lv := range rangeValues {
				// ask if price needs to be calculated given the ad price
				if _, ok := rangeValues["calculate"]; ok {
					adMap := ad.GetFieldsMapString()

					// if ad currency is 'peso', divide price by UF
					log.Printf("ad currency %v", adMap["currency"])
					var adPrice float64
					uf, _ := interactor.IndicatorsRepository.GetUF()
					if adMap["currency"] == "peso" {
						adPrice, _ = strconv.ParseFloat(adMap["price"], 64)
						adPrice /= uf

						log.Printf("ad price uf value %v", adPrice)
					}

					minusPrice, _ := strconv.Atoi(rangeValues["gte"].(string))
					plusPrice, _ := strconv.Atoi(rangeValues["lte"].(string))
					minPrice := adPrice - float64(minusPrice)
					maxPrice := adPrice + float64(plusPrice)

					log.Printf("min %v max %v ad %v ad %v", minPrice, maxPrice, adPrice, adMap["price"])

					rng["gte"] = fmt.Sprintf("%v", minPrice)
					rng["lte"] = fmt.Sprintf("%v", maxPrice)
					rng["type"] = rangeValues["type"].(string)
					rng["uf"] = fmt.Sprintf("%v", uf)
					out[rangeKey] = rng
				} else {
					rng[lk] = lv.(string)
					log.Printf("lk %s lv %s", lk, lv)
					out[rangeKey] = rng
				}
			}
		}
	}
	log.Printf("out %+v", out)
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
