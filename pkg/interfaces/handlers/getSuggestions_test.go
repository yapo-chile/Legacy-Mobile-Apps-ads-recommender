package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/domain"
)

type mockGetSuggestions struct {
	mock.Mock
}

func (m *mockGetSuggestions) GetSuggestions(
	listID string,
	optionalParams []string,
	size, from int,
	carouselType string,
) (ads []domain.Ad, err error) {
	args := m.Called(listID, size, from)
	return args.Get(0).([]domain.Ad), args.Error(1)
}

type mockDataMapping struct {
	mock.Mock
}

func (m *mockDataMapping) Get(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func TestGetSuggestionsHandlerInput(t *testing.T) {
	m := mockGetSuggestions{}
	mMockInputRequest := MockInputRequest{}
	mMockTargetRequest := MockTargetRequest{}
	mMockInputRequest.On(
		"Set", mock.AnythingOfType("*handlers.getSuggestionsHandlerInput"),
	).Return(&mMockTargetRequest)
	mMockTargetRequest.On("FromPath").Return()
	mMockTargetRequest.On("FromQuery").Return()

	h := GetSuggestionsHandler{
		Interactor: &m,
	}
	input := h.Input(&mMockInputRequest)

	var expected *getSuggestionsHandlerInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
	mMockTargetRequest.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
}

func TestGetSuggestionsHandlerResponseIg(t *testing.T) {
	response := &goutils.Response{Code: http.StatusOK}
	input := &getSuggestionsHandlerInput{}
	getter := MakeMockInputGetter(input, response)
	h := GetSuggestionsHandler{}
	r := h.Execute(getter)
	expected := response
	assert.Equal(t, expected, r)
}
func TestGetSuggestionsHandlerOK(t *testing.T) {
	mInteractor := &mockGetSuggestions{}
	timeT, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:10:10")
	ad := domain.Ad{
		ListID:   1,
		ListTime: timeT,
	}
	mInteractor.On(
		"GetSuggestions", mock.Anything, mock.Anything, mock.Anything,
	).Return([]domain.Ad{ad}, nil)
	h := GetSuggestionsHandler{
		Interactor: mInteractor,
	}
	input := &getSuggestionsHandlerInput{
		ListID: "1",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: getSuggestionsHandlerOutput{
			Ads: []AdsOutput{
				{
					ListID: "1",
					Date:   "2020-01-01 10:10:10",
				},
			},
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestGetProSuggestionsHandlerError(t *testing.T) {
	mInteractor := &mockGetSuggestions{}
	err := fmt.Errorf("err")
	mInteractor.On(
		"GetSuggestions", mock.Anything, mock.Anything, mock.Anything,
	).Return([]domain.Ad{}, err)

	h := GetSuggestionsHandler{
		Interactor: mInteractor,
	}
	input := &getSuggestionsHandlerInput{
		ListID: "1",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusInternalServerError,
		Body: &goutils.GenericError{
			ErrorMessage: err.Error(),
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestGetProSuggestionsHandlerEmptyResult(t *testing.T) {
	mInteractor := &mockGetSuggestions{}
	mInteractor.On(
		"GetSuggestions", mock.Anything, mock.Anything, mock.Anything,
	).Return([]domain.Ad{}, nil)
	h := GetSuggestionsHandler{
		Interactor: mInteractor,
	}
	input := &getSuggestionsHandlerInput{
		ListID: "1",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestGetProSuggestionsHandlerUFCurrency(t *testing.T) {
	timeT, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:10:10")
	ads := []domain.Ad{
		{
			ListID:   1,
			Currency: "uf",
			Price:    1000000,
			ListTime: timeT,
		},
	}
	h := GetSuggestionsHandler{
		CurrencySymbol:      "$",
		UnitOfAccountSymbol: "UF",
	}
	r := h.setOutput(ads, []string{})

	expected := getSuggestionsHandlerOutput{
		Ads: []AdsOutput{
			{
				ListID:   "1",
				Currency: "UF",
				Price:    10000,
				Date:     "2020-01-01 10:10:10",
			},
		},
	}
	assert.Equal(t, expected, r)
}
func TestGetProSuggestionsHandlerOtherCurrency(t *testing.T) {
	timeT, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:10:10")

	ads := []domain.Ad{
		{
			ListID:   1,
			Currency: "peso",
			Price:    1000000,
			ListTime: timeT,
		},
	}
	h := GetSuggestionsHandler{
		CurrencySymbol:      "$",
		UnitOfAccountSymbol: "UF",
	}
	r := h.setOutput(ads, []string{})

	expected := getSuggestionsHandlerOutput{
		Ads: []AdsOutput{
			{
				ListID:   "1",
				Currency: "$",
				Price:    1000000,
				Date:     "2020-01-01 10:10:10",
			},
		},
	}
	assert.Equal(t, expected, r)
}
func TestGetProSuggestionsHandlerOptionalParams(t *testing.T) {
	mRegions := mockDataMapping{}
	mCategories := mockDataMapping{}
	timeT, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:10:10")
	ads := []domain.Ad{
		{
			ListID:        1,
			CategoryID:    2020,
			Currency:      "peso",
			Price:         1000000,
			PublisherType: domain.Pro,
			ListTime:      timeT,
			CommuneID:     250,
			Type:          "sell",
			AdParams: map[string]string{
				"brand":   "TOYOTA",
				"Mileage": "1234"},
		},
	}
	mRegions.On("Get", mock.AnythingOfType("string")).Return("Metropolitana").Once()
	mCategories.On("Get", mock.AnythingOfType("string")).Return("autos, camionetas y 4x4").Once()
	mCategories.On("Get", mock.AnythingOfType("string")).Return("vehiculos").Once()
	h := GetSuggestionsHandler{
		CurrencySymbol:      "$",
		UnitOfAccountSymbol: "UF",
		Regions:             &mRegions,
		Categories:          &mCategories,
	}
	r := h.setOutput(ads, []string{
		"publisherType",
		"category",
		"communes",
		"region",
		"mileage",
		"type",
		"brand",
		"regDate",
		"unknown"})

	expected := getSuggestionsHandlerOutput{
		Ads: []AdsOutput{
			{
				ListID:              "1",
				Category:            "2020",
				Currency:            "$",
				Price:               1000000,
				PublisherType:       "pro",
				Brand:               "TOYOTA",
				Mileage:             "1234",
				Region:              "0",
				RegionDescription:   "Metropolitana",
				Communes:            "250",
				Date:                "2020-01-01 10:10:10",
				Type:                "s",
				CategoryDescription: "vehiculos > autos, camionetas y 4x4",
			},
		},
	}
	assert.Equal(t, expected, r)
	mRegions.AssertExpectations(t)
	mCategories.AssertExpectations(t)
}

func TestGetProSuggestionsSetWithoutSubcategory(t *testing.T) {
	mCategories := mockDataMapping{}
	timeT, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:10:10")

	ads := []domain.Ad{
		{
			ListID:     1,
			CategoryID: 2000,
			Currency:   "peso",
			Price:      1000000,
			ListTime:   timeT,
		},
	}
	mCategories.On("Get", mock.AnythingOfType("string")).Return("vehiculos").Once()
	h := GetSuggestionsHandler{
		CurrencySymbol:      "$",
		UnitOfAccountSymbol: "UF",
		Categories:          &mCategories,
	}
	r := h.setOutput(ads, []string{"category"})

	expected := getSuggestionsHandlerOutput{
		Ads: []AdsOutput{
			{
				ListID:              "1",
				Category:            "2000",
				Currency:            "$",
				Price:               1000000,
				Date:                "2020-01-01 10:10:10",
				CategoryDescription: "vehiculos",
			},
		},
	}
	assert.Equal(t, expected, r)
	mCategories.AssertExpectations(t)
}
func TestAddOptionalParamFalse(t *testing.T) {
	output := AdsOutput{}
	result := output.addOptionalParam("test", "a")
	expected := false
	assert.Equal(t, expected, result)
}

func TestSetFieldFalse(t *testing.T) {
	output := AdsOutput{}
	result := output.setField("test", "a")
	expected := false
	assert.Equal(t, expected, result)
}

func TestSetFixedURL(t *testing.T) {
	expected := "example.html"
	result := fixedURL("example")
	assert.Equal(t, expected, result)
}
