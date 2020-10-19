package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/Yapo/goutils"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/usecases"
)

// DataMapping allows get specific configuration params from etcd
type DataMapping interface {
	Get(string) string
}

// GetProSuggestionsHandler implements the handler interface and responds to
type GetProSuggestionsHandler struct {
	Interactor          usecases.GetSuggestionsInteractor
	CurrencySymbol      string
	UnitOfAccountSymbol string
	Regions             DataMapping
	Categories          DataMapping
}

type getProSuggestionsHandlerInput struct {
	ListID         string   `path:"listID"`
	Limit          int      `query:"limit"`
	OptionalParams []string `query:"params"`
}

// getProSuggestionsHandlerOutput struct that represents presenter output.
// This is the schema of endpoint response
type getProSuggestionsHandlerOutput struct {
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
func (*GetProSuggestionsHandler) Input(ir InputRequest) HandlerInput {
	input := getProSuggestionsHandlerInput{}
	ir.Set(&input).FromPath().FromQuery()
	return &input
}

// Execute is the main function of the GetProSuggestions handler
func (h *GetProSuggestionsHandler) Execute(ig InputGetter) *goutils.Response {
	input, err := ig()
	if err != nil {
		return err
	}
	in := input.(*getProSuggestionsHandlerInput)
	results, errSuggestions := h.Interactor.GetProSuggestions(in.ListID, in.OptionalParams, in.Limit, 0)
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
func (h *GetProSuggestionsHandler) setOutput(
	ads []domain.Ad, optionalParams []string,
) (out getProSuggestionsHandlerOutput) {
	for _, ad := range ads {
		adOutTemp := AdsOutput{
			ListID: strconv.FormatInt(ad.ListID, 10),
			Title:  ad.Subject,
			Price:  ad.Price,
			Date:   ad.ListTime.Format("2006-01-02 15:04:05"),
			Image: imageOutput{
				Full:   ad.Image.Full,
				Medium: ad.Image.Medium,
				Small:  ad.Image.Small,
			},
			URL: ad.URL,
		}
		if ad.Currency == "uf" {
			adOutTemp.Currency = h.UnitOfAccountSymbol
			adOutTemp.Price /= 100
		} else {
			adOutTemp.Currency = h.CurrencySymbol
		}
		// set optional params
		params := ad.GetAllFields()
		for _, param := range optionalParams {
			if val, ok := params[strings.ToLower(param)]; ok {
				if value := fmt.Sprintf("%v", val); value != "" {
					adOutTemp.addOptionalParam(param, value)
				}
			}
			if strings.ToLower(param) == "publishertype" {
				adOutTemp.PublisherType = string(ad.PublisherType)
			}
			if strings.ToLower(param) == "region" {
				adOutTemp.Region = strconv.FormatInt(ad.RegionID, 10)
				adOutTemp.SetRegion(h.Regions)
			}
			if strings.ToLower(param) == "communes" {
				adOutTemp.Communes = strconv.FormatInt(ad.CommuneID, 10)
				adOutTemp.CommunesDescription = ad.Commune
			}
			if strings.ToLower(param) == "type" {
				if adOutTemp.Type != "" {
					adOutTemp.Type = strings.ToLower(string(adOutTemp.Type[0]))
				}
			}
			if strings.ToLower(param) == "category" {
				adOutTemp.Category = strconv.FormatInt(ad.CategoryID, 10)
				adOutTemp.SetCategory(h.Categories)
			}
		}
		out.Ads = append(out.Ads, adOutTemp)
	}
	return out
}
