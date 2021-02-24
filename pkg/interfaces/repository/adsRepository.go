package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/usecases"
)

var notAlphaNumbericRegex = regexp.MustCompile("[^a-zA-Z0-9]+")
var specialCases = strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o",
	"ú", "u", "'", "", "ñ", "n")

// adsRepository contains the required variables and functions
// to call the methods on elastic handler
type adsRepository struct {
	elasticHandler  ElasticSearchHandler
	regionsConf     DataMapping
	imageServerLink string
	index           string
	resultSize      int
	from            int
	queryTemplates  map[string]*template.Template
}

// Hit represent a query match on elasticsearch
type Hit struct {
	Source usecases.Ad `json:"_source"`
}

// Hits is a slice of Hits on elasticsearch
type Hits []Hit

// HitsParent is the root response of a elastic query hit
type HitsParent struct {
	Hits Hits `json:"hits"`
}

// elasticResponse is the response that elastic search gives to a query
type elasticResponse struct {
	HitsParent HitsParent `json:"hits"`
}

// NewAdsRepository return a new ads repositoryinstance
func NewAdsRepository(
	handler ElasticSearchHandler,
	regionsConf DataMapping,
	queryTemplates map[string]*template.Template,
	imageServerLink, index string, resultSize, from int) usecases.AdsRepository {
	return &adsRepository{
		elasticHandler:  handler,
		imageServerLink: imageServerLink,
		index:           index,
		resultSize:      resultSize,
		from:            from,
		queryTemplates:  queryTemplates,
		regionsConf:     regionsConf,
	}
}

// GetAd returns a unique Ad object using listID
func (repo *adsRepository) GetAd(listID string) (ad domain.Ad, err error) {
	params := map[string]string{
		"ListID": listID,
	}
	ads, err := repo.getAdsProcess("getAd", params, 0, 0)
	if err != nil {
		return
	}
	if len(ads) != 1 {
		err = fmt.Errorf("get ad fails to get it, len: %d", len(ads))
		return
	}
	return ads[0], nil
}

// GetAds returns a slice of Ad object using mandatory parameters (musts),
// optional parameters(shoulds), exclude results if param is on ad(mustsNot)
// and aditional keyword filters (filters) to get ads related to this terms.
func (repo *adsRepository) GetAds(
	musts, shoulds, mustsNot, filters, priceRange, decay map[string]string,
	queryString []map[string]string,
	size, from int,
) (ads []domain.Ad, err error) {

	mustsParams := repo.getBoolParameters(musts)
	mustsNotParams := repo.getBoolParameters(mustsNot)
	shouldsParams := repo.getBoolParameters(shoulds)
	filtersParams := repo.getFilters(filters)
	queryStringParams := repo.getQueryString(queryString)

	if len(priceRange) > 0 {
		priceParams := repo.processPriceTemplate(priceRange)
		switch priceRange["type"] {
		case "must":
			mustsParams = joinParams(mustsParams, priceParams)
		case "mustNot":
			mustsNotParams = joinParams(mustsNotParams, priceParams)
		case "should":
			shouldsParams = joinParams(shouldsParams, priceParams)
		case "filter":
			filtersParams = joinParams(filtersParams, priceParams)
		}
	}

	if len(queryStringParams) > 0 {
		mustsParams = joinParams(mustsParams, queryStringParams)
	}

	params := map[string]string{
		"Musts":    mustsParams,
		"MustsNot": mustsNotParams,
		"Shoulds":  shouldsParams,
		"Filters":  filtersParams,
		"Name":     decay["name"],
		"Field":    decay["field"],
		"Origin":   decay["origin"],
		"Offset":   decay["offset"],
		"Scale":    decay["scale"],
	}
	return repo.getAdsProcess("getAds", params, size, from)
}

// getAdsProcess executes a query to elastic search through the elastic handler
// and process the response. It returns an ads slice
func (repo *adsRepository) getAdsProcess(
	templateName string,
	params map[string]string,
	size, from int,
) (ads []domain.Ad, err error) {
	query, err := repo.ProcessTemplate(templateName, params)
	log.Printf("query: %+v", query)
	if err != nil {
		return
	}
	if size == 0 {
		size = repo.resultSize
	}
	if from == 0 {
		from = repo.from
	}
	response, err := repo.elasticHandler.Search(repo.index, query, size, from)
	if err != nil {
		return
	}
	var parsed elasticResponse

	if err = json.Unmarshal([]byte(response), &parsed); err != nil {
		return
	}
	for _, hit := range parsed.HitsParent.Hits {
		ads = append(ads, repo.fillAd(hit.Source))
	}
	return ads, nil
}

// getBoolParameters returns a string with bool parameters
// to be used on a query as must, should or must_not
func (repo *adsRepository) getBoolParameters(params map[string]string) string {
	return repo.getParams(params, `{"match": {"%s": "%s"}}`)
}

// getBoolParameters returns a string with bool parameters
// to be used on a query as must, should or must_not
func (repo *adsRepository) getQueryString(params []map[string]string) string {
	var out string
	for _, param := range params {
		o := repo.getParams(param, `"%s": "%s"`)
		out = joinParams(out, fmt.Sprintf(`{"query_string": {%s}}`, o))
	}
	return out
}

// processPriceTemplate returns the range query template as string
// to be used in the final query
func (repo *adsRepository) processPriceTemplate(priceRange map[string]string) string {
	params := map[string]string{
		"PriceMin": priceRange["gte"],
		"PriceMax": priceRange["lte"],
		"UF":       priceRange["uf"],
	}
	query, err := repo.ProcessTemplate("priceScript", params)
	if err != nil {
		return ""
	}
	return query
}

// getFilters returns a string with filters to be used on a query
func (repo *adsRepository) getFilters(filters map[string]string) string {
	return repo.getParams(filters, `{"term": {"%s.keyword": "%s"}}`)
}

// getParams returns a string to be used on a query
func (repo *adsRepository) getParams(params map[string]string, condition string) string {
	var paramsStr strings.Builder
	if len(params) > 0 {
		keys := sortedKeys(params)
		for i, k := range keys {
			if i > 0 {
				paramsStr.WriteString(`, `)
			}
			paramsStr.WriteString(fmt.Sprintf(condition, k, params[k]))
		}
	}
	return paramsStr.String()
}

// ProcessTemplate process the query data and returns a template as string.
// If something goes wrong returns empty string and error
func (repo *adsRepository) ProcessTemplate(template string, params map[string]string) (string, error) {
	if val, ok := repo.queryTemplates[template]; ok {
		var processedTemplate bytes.Buffer
		if err := val.Execute(&processedTemplate, params); err != nil {
			return "", err
		}
		return processedTemplate.String(), nil
	}
	return "", fmt.Errorf("template not found")
}

// getMainImage gets the main image for required ad using media struct
func (repo *adsRepository) getMainImage(imgs []usecases.AdMedia) domain.Image {
	if len(imgs) == 0 {
		return domain.Image{}
	}
	for _, img := range imgs {
		if img.SeqNo == 0 {
			return repo.fillImage(img.ID)
		}
	}
	return repo.fillImage(imgs[0].ID)
}

// fillAd parse data from Ad struct on usecases to Ad domain object
func (repo *adsRepository) fillAd(ad usecases.Ad) domain.Ad {
	return domain.Ad{
		ListID:        ad.ListID,
		UserID:        ad.UserID,
		CategoryID:    ad.CategoryID,
		Category:      ad.Category,
		Type:          ad.Type,
		CommuneID:     ad.CommuneID,
		RegionID:      ad.RegionID,
		Phone:         ad.Phone,
		Region:        ad.Region,
		Commune:       ad.Commune,
		SubCategory:   ad.SubCategory,
		Name:          ad.Name,
		Body:          ad.Body,
		OldPrice:      ad.OldPrice,
		ListTime:      ad.ListTime,
		Subject:       ad.Subject,
		Price:         ad.Price,
		PublisherType: ad.PublisherType,
		Currency:      ad.Params["Currency"],
		URL:           repo.fillURL(ad.Subject, ad.ListID, ad.RegionID),
		Image:         repo.getMainImage(ad.Media),
		AdParams:      ad.Params,
	}
}

// fillImage parses the image id to domain Image struct
func (repo *adsRepository) fillImage(id int) domain.Image {
	IDstr := fmt.Sprintf("%010d", id)
	return domain.Image{
		Full:   fmt.Sprintf(repo.imageServerLink, "images", IDstr[:2], IDstr),
		Medium: fmt.Sprintf(repo.imageServerLink, "thumbsli", IDstr[:2], IDstr),
		Small:  fmt.Sprintf(repo.imageServerLink, "thumbs", IDstr[:2], IDstr),
	}
}

// fillURL returns the main URL to visit an ad on site
func (repo *adsRepository) fillURL(subject string, listID, regionID int64) string {
	regionKey := fmt.Sprintf("region.%d.link", regionID)
	regionName := repo.regionsConf.Get(regionKey)
	return "/" + strings.Join(
		[]string{
			notAlphaNumbericRegex.ReplaceAllString(
				specialCases.Replace(strings.ToLower(regionName)), "_"),
			notAlphaNumbericRegex.ReplaceAllString(
				specialCases.Replace(strings.ToLower(subject)), "_") +
				"_" + strconv.FormatInt(listID, 10),
		},
		"/",
	)
}

func sortedKeys(m map[string]string) (keys []string) {
	keys = make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// joinParams concatenates the strings in an array with a , and
// returns the resulting string
func joinParams(params ...string) (output string) {
	paramsSlice := make([]string, 0)
	for _, val := range params {
		if len(val) > 0 {
			paramsSlice = append(paramsSlice, val)
		}
	}
	return strings.Join(paramsSlice, ",")
}
