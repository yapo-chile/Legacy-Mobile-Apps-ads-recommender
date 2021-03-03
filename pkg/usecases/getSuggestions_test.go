package usecases

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
)

type mockGetSuggestionsLogger struct {
	mock.Mock
}

func (m *mockGetSuggestionsLogger) LimitExceeded(size, maxDisplayedAds, defaultAdsQty int) {
	m.Called(size, maxDisplayedAds, defaultAdsQty)
}
func (m *mockGetSuggestionsLogger) MinimumQtyNotEnough(size, minDisplayedAds, defaultAdsQty int) {
	m.Called(size, minDisplayedAds, defaultAdsQty)
}
func (m *mockGetSuggestionsLogger) ErrorGettingAd(listID string, err error) {
	m.Called(listID, err)
}
func (m *mockGetSuggestionsLogger) ErrorGettingAds(musts, shoulds, mustsNot map[string]string, err error) {
	m.Called(musts, shoulds, mustsNot, err)
}
func (m *mockGetSuggestionsLogger) NotEnoughAds(listID string, lenAds int) {
	m.Called(listID, lenAds)
}
func (m *mockGetSuggestionsLogger) ErrorGettingAdsContact(listID string, err error) {
	m.Called(listID, err)
}
func (m *mockGetSuggestionsLogger) ErrorGettingUF(err error) {
	m.Called(err)
}
func (m *mockGetSuggestionsLogger) InvalidCarousel(carousel string) {
	m.Called(carousel)
}

type mockAdsRepository struct {
	mock.Mock
}

func (m *mockAdsRepository) GetAd(listID string) (domain.Ad, error) {
	args := m.Called(listID)
	return args.Get(0).(domain.Ad), args.Error(1)
}
func (m *mockAdsRepository) GetAds(
	musts, shoulds, mustsNot, filters, priceRange, decay map[string]string,
	queryStrings []map[string]string,
	size, from int,
) ([]domain.Ad, error) {
	args := m.Called(musts, shoulds, mustsNot, filters, priceRange, decay, queryStrings, size, from)
	return args.Get(0).([]domain.Ad), args.Error(1)
}

type mockAdContactRepository struct {
	mock.Mock
}

func (m *mockAdContactRepository) GetAdsPhone(
	suggestions []domain.Ad,
) (adsResult map[string]string, err error) {
	args := m.Called(suggestions)
	return args.Get(0).(map[string]string), args.Error(1)
}

type mockIndicatorsRepository struct {
	mock.Mock
}

func (m *mockIndicatorsRepository) GetUF() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func getDefaultSuggestionParams() (out map[string]map[string][]interface{}) {
	out = make(map[string]map[string][]interface{})
	out["default"] = make(map[string][]interface{})

	out["default"]["must"] = []interface{}{"Category", "SubCategory"}
	out["default"]["should"] = []interface{}{
		"Params.BrandID",
		"Params.ModelID",
		"Params.Regdate",
		"Params.Brand",
		"Params.Model",
	}
	out["default"]["mustNot"] = []interface{}{"ListID"}
	out["default"]["queryString"] = []interface{}{
		map[string]interface{}{
			"query":        "(pro OR professional)",
			"defaultField": "PublisherType",
		},
	}
	out["default"]["decayFunc"] = []interface{}{
		map[string]interface{}{
			"name":   "gauss",
			"field":  "ListTime",
			"origin": "now/1d",
			"offset": "1d",
			"scale":  "7d",
		},
	}
	return
}

func getSuggestionParams(
	carousel string,
	params ...map[string][]interface{}) (out map[string]map[string][]interface{}) {
	out = make(map[string]map[string][]interface{})
	out["default"] = getDefaultSuggestionParams()["default"]
	if carousel != "default" {
		out[carousel] = make(map[string][]interface{})
	}

	for _, param := range params {
		for key, val := range param {
			out[carousel][key] = val
		}
	}
	return
}

func TestGetProSuggestionsOK(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mIndicatorsRepo := mockIndicatorsRepository{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	ads := []domain.Ad{{ListID: 2, Category: "test"}}
	priceRange := map[string][]interface{}{
		"priceRange": {
			map[string]interface{}{
				"gte": "100",
				"lte": "200",
			},
		},
	}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(ads, nil)
	mIndicatorsRepo.On("GetUF").Return(float64(28000), nil)
	i := GetSuggestions{
		SuggestionsRepo:      &mAdsRepo,
		IndicatorsRepository: &mIndicatorsRepo,
		MinDisplayedAds:      1,
		MaxDisplayedAds:      1,
		SuggestionsParams:    getSuggestionParams("default", priceRange),
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{}, 1, 0, "default")
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mIndicatorsRepo.AssertExpectations(t)
}

func TestGetProSuggestionsMaxDisplayedAds(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mLogger := mockGetSuggestionsLogger{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	ads := []domain.Ad{{ListID: 2, Category: "test"}}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(ads, nil)
	mLogger.On("LimitExceeded", mock.Anything, mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo:   &mAdsRepo,
		MinDisplayedAds:   1,
		MaxDisplayedAds:   1,
		RequestedAdsQty:   2,
		Logger:            &mLogger,
		SuggestionsParams: getSuggestionParams("default"),
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{}, 2, 0, "default")
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetProSuggestionsMinDisplayedAds(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mLogger := mockGetSuggestionsLogger{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	ads := []domain.Ad{{ListID: 2, Category: "test"}, {ListID: 3, Category: "test"}}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(ads, nil)
	mLogger.On("MinimumQtyNotEnough", mock.Anything, mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo:   &mAdsRepo,
		SuggestionsParams: getSuggestionParams("default"),
		MinDisplayedAds:   2,
		MaxDisplayedAds:   2,
		RequestedAdsQty:   2,
		Logger:            &mLogger,
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{}, 1, 0, "default")
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetProSuggestionsNotEnoughAds(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mLogger := mockGetSuggestionsLogger{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	ads := []domain.Ad{{ListID: 2, Category: "test"}}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(ads, nil)
	mLogger.On("NotEnoughAds", mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo:   &mAdsRepo,
		SuggestionsParams: getSuggestionParams("default"),
		MinDisplayedAds:   2,
		MaxDisplayedAds:   2,
		RequestedAdsQty:   2,
		Logger:            &mLogger,
	}
	expected := []domain.Ad{}
	output, err := i.GetProSuggestions("1", []string{}, 2, 0, "default")
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetProSuggestionsGetAdErr(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mLogger := mockGetSuggestionsLogger{}
	mAdsRepo.On("GetAd", mock.Anything).Return(domain.Ad{}, fmt.Errorf("error"))
	mLogger.On("ErrorGettingAd", mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo:   &mAdsRepo,
		SuggestionsParams: getSuggestionParams("default"),
		MinDisplayedAds:   1,
		MaxDisplayedAds:   2,
		Logger:            &mLogger,
	}
	var expected []domain.Ad
	output, err := i.GetProSuggestions("1", []string{}, 1, 0, "default")
	assert.Error(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetProSuggestionsGetAdsErr(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mLogger := mockGetSuggestionsLogger{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return([]domain.Ad{}, fmt.Errorf(""))
	mLogger.On("ErrorGettingAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo:   &mAdsRepo,
		SuggestionsParams: getSuggestionParams("default"),
		MinDisplayedAds:   1,
		MaxDisplayedAds:   1,
		RequestedAdsQty:   2,
		Logger:            &mLogger,
	}
	expected := []domain.Ad{}
	output, err := i.GetProSuggestions("1", []string{}, 1, 0, "default")
	assert.Error(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetProSuggestionsOKWithPhoneLink(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mAdContactRepo := mockAdContactRepository{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	ads := []domain.Ad{{ListID: 2, Category: "test"}}
	phones := map[string]string{"2": "998765432"}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(ads, nil)
	mAdContactRepo.On("GetAdsPhone", mock.Anything).Return(phones, nil)
	i := GetSuggestions{
		SuggestionsRepo:   &mAdsRepo,
		SuggestionsParams: getSuggestionParams("default"),
		AdContact:         &mAdContactRepo,
		MinDisplayedAds:   1,
		MaxDisplayedAds:   2,
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{"phonelink"}, 1, 0, "default")
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mAdContactRepo.AssertExpectations(t)
}

func TestGetProSuggestionsWithPhoneLinkErr(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mAdContactRepo := mockAdContactRepository{}
	mLogger := mockGetSuggestionsLogger{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	ads := []domain.Ad{{ListID: 2, Category: "test"}}
	phones := map[string]string{"2": "998765432"}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(ads, nil)
	mAdContactRepo.On("GetAdsPhone", mock.Anything).Return(phones, fmt.Errorf("error"))
	mLogger.On("ErrorGettingAdsContact", mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo:   &mAdsRepo,
		AdContact:         &mAdContactRepo,
		SuggestionsParams: getSuggestionParams("default"),
		MinDisplayedAds:   1,
		MaxDisplayedAds:   2,
		Logger:            &mLogger,
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{"phonelink"}, 1, 0, "default")
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mAdContactRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetProSuggestionsGetAdsInvalidCarousel(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mLogger := mockGetSuggestionsLogger{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mLogger.On("InvalidCarousel", mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo:   &mAdsRepo,
		SuggestionsParams: getSuggestionParams("default"),
		MinDisplayedAds:   1,
		MaxDisplayedAds:   1,
		RequestedAdsQty:   2,
		Logger:            &mLogger,
	}
	expected := []domain.Ad{}
	output, err := i.GetProSuggestions("1", []string{}, 1, 0, "not_a_carousel")
	assert.Error(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetProSuggestionsGetAdsErrUF(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mIndicatorsRepo := mockIndicatorsRepository{}
	mLogger := mockGetSuggestionsLogger{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	priceRange := map[string][]interface{}{
		"priceRange": {
			map[string]interface{}{
				"gte": "100",
				"lte": "200",
			},
		},
	}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mIndicatorsRepo.On("GetUF").Return(float64(0), fmt.Errorf("error"))
	mLogger.On("ErrorGettingUF", mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo:      &mAdsRepo,
		IndicatorsRepository: &mIndicatorsRepo,
		SuggestionsParams:    getSuggestionParams("default", priceRange),
		MinDisplayedAds:      1,
		MaxDisplayedAds:      1,
		RequestedAdsQty:      2,
		Logger:               &mLogger,
	}
	expected := []domain.Ad{}
	output, err := i.GetProSuggestions("1", []string{}, 1, 0, "default")
	assert.Error(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mIndicatorsRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetPriceRangeOK(t *testing.T) {
	ad := domain.Ad{
		Price:    1000,
		Currency: "uf",
	}
	mIndicatorsRepo := mockIndicatorsRepository{}
	mIndicatorsRepo.On("GetUF").Return(float64(10), nil)
	i := GetSuggestions{IndicatorsRepository: &mIndicatorsRepo}
	testCases := []struct {
		name       string
		priceRange []interface{}
		expected   map[string]string
	}{
		{
			"gte and lte only",
			[]interface{}{map[string]interface{}{"gte": "1000", "lte": "2000"}},
			map[string]string{"gte": "1000", "lte": "2000", "type": "must", "uf": "10"},
		},
		{
			"calculate price",
			[]interface{}{map[string]interface{}{"gte": "10", "lte": "20", "calculate": "true"}},
			map[string]string{"gte": "0", "lte": "30", "type": "must", "uf": "10"},
		},
		{
			"should type",
			[]interface{}{map[string]interface{}{"gte": "1000", "lte": "2000", "type": "should"}},
			map[string]string{"gte": "1000", "lte": "2000", "type": "should", "uf": "10"},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			output, err := i.getPriceRange(ad, tc.priceRange)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, output)
			mIndicatorsRepo.AssertExpectations(t)
		})
	}
}

func TestCalculateMinMaxPriceRange(t *testing.T) {
	testCases := []struct {
		name                     string
		adPrice, uf              float64
		adCurrency               string
		minusValue, plusValue    int
		expectedMin, expectedMax string
	}{
		{
			"integer peso currency",
			1000, 10, "peso", 10, 10, "90", "110",
		},
		{
			"float peso currency",
			1055, 10.55, "peso", 10, 10, "90", "110",
		},
		{
			"integer uf currency",
			1000, 10, "uf", 10, 10, "0", "20",
		},
		{
			"float uf currency",
			800000, 10, "uf", 1000, 1000, "7000", "9000",
		},
		{
			"negative min value uf",
			80000, 10, "uf", 1000, 1000, "-200", "1800",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			min, max := calculateMinMaxPriceRange(tc.adPrice, tc.uf, tc.adCurrency, tc.minusValue, tc.plusValue)
			assert.Equal(t, min, tc.expectedMin)
			assert.Equal(t, max, tc.expectedMax)
		})
	}
}

func TestGetSliceParams(t *testing.T) {
	testCases := []struct {
		name              string
		adMap             map[string]string
		suggestionsParams []interface{}
		expected          map[string]string
	}{
		{
			"simple params",
			map[string]string{"a": "testA", "b": "testB"},
			[]interface{}{"a", "b"},
			map[string]string{"a": "testA", "b": "testB"},
		},
		{
			"compound params",
			map[string]string{"a": "testA", "b": "testB"},
			[]interface{}{"Params.a", "Params.b"},
			map[string]string{"Params.a": "testA", "Params.b": "testB"},
		},
		{
			"mixed params",
			map[string]string{"a": "testA", "b": "testB"},
			[]interface{}{"Params.a", "b"},
			map[string]string{"Params.a": "testA", "b": "testB"},
		},
		{
			"empty params",
			map[string]string{"a": "testA", "b": "testB"},
			[]interface{}{},
			map[string]string{},
		},
		{
			"ad with empty params",
			map[string]string{},
			[]interface{}{"Params.a", "b"},
			map[string]string{},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			params := getSliceParams(tc.adMap, tc.suggestionsParams)
			assert.Equal(t, tc.expected, params)
		})
	}
}

func TestGetDecayFunctionParams(t *testing.T) {
	decayFunc := map[string][]interface{}{
		"decayFunc": {
			map[string]interface{}{
				"name":   "gauss",
				"field":  "ListTime",
				"origin": "now/5d",
				"offset": "5d",
				"scale":  "69d",
			},
		},
	}
	validCarousel := getSuggestionParams("valid", decayFunc)
	defaultCarousel := getSuggestionParams("notfound")

	testCases := []struct {
		name          string
		decayFuncConf map[string]map[string][]interface{}
		carouselType  string
		expected      map[string]string
	}{
		{
			"carousel with decay func",
			validCarousel, "valid",
			map[string]string{
				"name":   "gauss",
				"field":  "ListTime",
				"origin": "now/5d",
				"offset": "5d",
				"scale":  "69d",
			},
		},
		{
			"carousel not found, default used",
			defaultCarousel, "notfound",
			map[string]string{
				"name":   "gauss",
				"field":  "ListTime",
				"origin": "now/1d",
				"offset": "1d",
				"scale":  "7d",
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			decayParams := getDecayFunctionParams(tc.decayFuncConf, tc.carouselType)
			assert.Equal(t, tc.expected, decayParams)
		})
	}
}

func TestGetQueryStringParams(t *testing.T) {
	testCases := []struct {
		name             string
		queryStringSlice []interface{}
		expected         []map[string]string
	}{
		{
			"one query param",
			[]interface{}{
				map[string]interface{}{
					"query":        "(pro OR professional)",
					"defaultField": "PublisherType",
				},
			},
			[]map[string]string{
				{
					"query":         "(pro OR professional)",
					"default_field": "PublisherType",
				},
			},
		},
		{
			"multiple query params",
			[]interface{}{
				map[string]interface{}{
					"query":        "(pro OR professional)",
					"defaultField": "PublisherType",
				},
				map[string]interface{}{
					"query":        "(sell OR buy)",
					"defaultField": "Type",
				},
			},
			[]map[string]string{
				{
					"query":         "(pro OR professional)",
					"default_field": "PublisherType",
				},
				{
					"query":         "(sell OR buy)",
					"default_field": "Type",
				},
			},
		},
		{
			"no query params",
			[]interface{}{},
			[]map[string]string{},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			decayParams := getQueryStringParams(tc.queryStringSlice)
			assert.Equal(t, tc.expected, decayParams)
		})
	}
}
