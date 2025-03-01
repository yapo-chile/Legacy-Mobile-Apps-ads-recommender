package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/Yapo/goutils"
	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/domain"
	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/usecases"
)

// DataMapping allows get specific configuration params from etcd
type DataMapping interface {
	Get(string) string
}

// GetSuggestionsHandler implements the handler interface and responds to
type GetSuggestionsHandler struct {
	Interactor          usecases.GetSuggestionsInteractor
	CurrencySymbol      string
	UnitOfAccountSymbol string
	Regions             DataMapping
	Categories          DataMapping
}

type getSuggestionsHandlerInput struct {
	ListID         string   `path:"listID"`
	From           int      `query:"from"`
	Limit          int      `query:"limit"`
	OptionalParams []string `query:"params"`
	CarouselType   string   `path:"carousel"`
}

// getProSuggestionsHandlerOutput struct that represents presenter output.
// This is the schema of endpoint response
type getSuggestionsHandlerOutput struct {
	Ads []AdsOutput `json:"ads"`
}

// AdsOutput struct that represents Ads schema output
type AdsOutput struct {
	ListID              string      `json:"id"`
	Title               string      `json:"title"`
	Price               float64     `json:"price"`
	Currency            string      `json:"currency"`
	Image               imageOutput `json:"images"`
	URL                 string      `json:"url"`
	Region              string      `json:"region,omitempty"`
	RegionDescription   string      `json:"regionDescription,omitempty"`
	Communes            string      `json:"communes,omitempty"`
	CommunesDescription string      `json:"communesDescription,omitempty"`
	Date                string      `json:"date,omitempty"`
	Description         string      `json:"body,omitempty"`
	Category            string      `json:"category,omitempty"`
	CategoryDescription string      `json:"categoryDescription,omitempty"`
	Brand               string      `json:"brand,omitempty"`
	BuiltYear           string      `json:"builtYear,omitempty"`
	Capacity            string      `json:"capacity,omitempty"`
	Cartype             string      `json:"cartype,omitempty"`
	City                string      `json:"city,omitempty"`
	Condition           string      `json:"condition,omitempty"`
	Cubiccms            string      `json:"cubiccms,omitempty"`
	Equipment           string      `json:"equipment,omitempty"`
	EstateType          string      `json:"estateType,omitempty"`
	Fuel                string      `json:"fuel,omitempty"`
	GearBox             string      `json:"gearbox,omitempty"`
	Mileage             string      `json:"mileage,omitempty"`
	Model               string      `json:"model,omitempty"`
	Regdate             string      `json:"regdate,omitempty"`
	Type                string      `json:"type,omitempty"`
	PhoneLink           string      `json:"phonelink,omitempty"`
	PublisherType       string      `json:"publisherType,omitempty"`
	BrandID             string      `json:"brandid,omitempty"`
	ModelID             string      `json:"modelid,omitempty"`
}

// Image struct that defines the internal structure of the images
// that related Ads endpoint will return
type imageOutput struct {
	Full   string `json:"full,omitempty"`
	Medium string `json:"medium,omitempty"`
	Small  string `json:"small,omitempty"`
}

// addOptionalParam sets a value on AdsOutput if name is a tag on struct
// returns true when all is done successfully, otherwise false
func (output *AdsOutput) addOptionalParam(name, value string) bool {
	val := reflect.ValueOf(output).Elem()
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		if fieldtag := t.Tag.Get("json"); fieldtag != "" && fieldtag != "-" {
			if commaIdx := strings.Index(fieldtag, ","); commaIdx > 0 {
				fieldtag = fieldtag[:commaIdx]
			}
			if strings.EqualFold(fieldtag, name) {
				return output.setField(t.Name, value)
			}
		}
	}
	return false
}

// setField sets a field with a value. If its done
// successfully return true, otherwise false
func (output *AdsOutput) setField(name, value string) bool {
	rv := reflect.ValueOf(output).Elem()
	field := rv.FieldByName(name)
	if field.IsValid() && field.CanSet() {
		if field.Kind() == reflect.String {
			field.SetString(value)
			return true
		}
	}
	return false
}

// SetRegion gets region label on Config
func (output *AdsOutput) SetRegion(regions DataMapping) {
	if output.Region != "" {
		if regionID, err := strconv.ParseInt(output.Region, 10, 64); err == nil {
			output.RegionDescription = regions.Get(
				fmt.Sprintf(
					"region.%d.name",
					regionID,
				),
			)
		}
	}
}

// SetCategory it gets category and main category labels on DataMapping
// and sets on Category field
func (output *AdsOutput) SetCategory(categories DataMapping) {
	if output.Category != "" {
		if categoryID, err := strconv.ParseInt(output.Category, 10, 64); err == nil {
			categoryDescription := categories.Get(fmt.Sprintf("%d", categoryID))
			mainCategoryID := (categoryID / 1000) * 1000
			// check if category has a main category
			if categoryID != mainCategoryID {
				mainCategoryDescription := categories.Get(fmt.Sprintf("%d", mainCategoryID))
				output.CategoryDescription = fmt.Sprintf("%s > %s", mainCategoryDescription, categoryDescription)
			} else {
				output.CategoryDescription = categoryDescription
			}
		}
	}
}

// Input returns a fresh, empty instance of getProSuggestionsHandlerInput
func (*GetSuggestionsHandler) Input(ir InputRequest) HandlerInput {
	input := getSuggestionsHandlerInput{}
	ir.Set(&input).FromPath().FromQuery()
	return &input
}

// Execute is the main function of the GetProSuggestions handler
func (h *GetSuggestionsHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		if response.Code == http.StatusOK {
			return response
		}
	}
	in := input.(*getSuggestionsHandlerInput)
	results, errSuggestions := h.Interactor.GetSuggestions(
		in.ListID,
		in.OptionalParams,
		in.Limit,
		in.From,
		in.CarouselType,
	)
	if errSuggestions != nil {
		return &goutils.Response{
			Code: http.StatusInternalServerError,
			Body: &goutils.GenericError{
				ErrorMessage: errSuggestions.Error(),
			},
		}
	}
	if len(results) == 0 {
		return &goutils.Response{
			Code: http.StatusNoContent,
		}
	}
	return &goutils.Response{
		Code: http.StatusOK,
		Body: h.setOutput(results, in.OptionalParams),
	}
}

// setOutput sets presenter to format the output response for getSuggestions usecase
func (h *GetSuggestionsHandler) setOutput(
	ads []domain.Ad, optionalParams []string,
) (out getSuggestionsHandlerOutput) {
	for _, ad := range ads {
		// get a map with all params on ads as string
		params := ad.GetFieldsMapString()
		adOutTemp := AdsOutput{
			ListID: params["listid"],
			Title:  params["subject"],
			Price:  ad.Price,
			Date:   params["listtime"],
			Image: imageOutput{
				Full:   ad.Image.Full,
				Medium: ad.Image.Medium,
				Small:  ad.Image.Small,
			},
			URL: fixedURL(params["url"]),
		}
		if ad.Currency == "uf" {
			adOutTemp.Currency = h.UnitOfAccountSymbol
		} else {
			adOutTemp.Currency = h.CurrencySymbol
		}

		// set optional params
		for _, optionalParam := range optionalParams {
			optionalParam = strings.ToLower(optionalParam)
			if val, ok := params[optionalParam]; ok {
				if val != "" {
					adOutTemp.addOptionalParam(optionalParam, val)
				}
			}
			if optionalParam == "publishertype" {
				adOutTemp.PublisherType = string(ad.PublisherType)
			}
			if optionalParam == "region" {
				adOutTemp.Region = strconv.FormatInt(ad.RegionID, 10)
				adOutTemp.SetRegion(h.Regions)
			}
			if optionalParam == "communes" {
				adOutTemp.Communes = strconv.FormatInt(ad.CommuneID, 10)
				adOutTemp.CommunesDescription = ad.Commune
			}
			if optionalParam == "type" {
				if adOutTemp.Type != "" {
					adOutTemp.Type = strings.ToLower(string(adOutTemp.Type[0]))
				}
			}
			if optionalParam == "category" {
				adOutTemp.Category = strconv.FormatInt(ad.CategoryID, 10)
				adOutTemp.SetCategory(h.Categories)
			}
		}
		out.Ads = append(out.Ads, adOutTemp)
	}
	return out
}

// fixedURL returns a valid page to redirect
func fixedURL(url string) string {
	if url != "" && !strings.HasSuffix(url, ".html") {
		url += ".html"
	}
	return url
}
