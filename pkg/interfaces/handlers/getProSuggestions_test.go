package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
)

type mockGetSuggestions struct {
	mock.Mock
}

func (m *mockGetSuggestions) GetProSuggestions(listID string, size, from int) (ads []domain.Ad, err error) {
	args := m.Called(listID, size, from)
	return args.Get(0).([]domain.Ad), args.Error(1)
}

func TestGetProSuggestionsHandlerInput(t *testing.T) {
	m := mockGetSuggestions{}
	mMockInputRequest := MockInputRequest{}
	mMockTargetRequest := MockTargetRequest{}
	mMockInputRequest.On(
		"Set", mock.AnythingOfType("*handlers.getProSuggestionsHandlerInput"),
	).Return(&mMockTargetRequest)
	mMockTargetRequest.On("FromPath").Return()
	mMockTargetRequest.On("FromQuery").Return()

	h := GetProSuggestionsHandler{
		Interactor: &m,
	}
	input := h.Input(&mMockInputRequest)

	var expected *getProSuggestionsHandlerInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
	mMockTargetRequest.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
}

func TestGetProSuggestionsHandlerErrIg(t *testing.T) {
	response := &goutils.Response{}

	input := &getProSuggestionsHandlerInput{
		ListID: "1",
	}
	getter := MakeMockInputGetter(input, response)

	h := GetProSuggestionsHandler{}
	r := h.Execute(getter)

	expected := response
	assert.Equal(t, expected, r)
}

func TestGetProSuggestionsHandlerOK(t *testing.T) {
	mInteractor := &mockGetSuggestions{}
	ad := domain.Ad{ListID: 1, CategoryID: 2020}
	mInteractor.On(
		"GetProSuggestions", mock.Anything, mock.Anything, mock.Anything,
	).Return([]domain.Ad{ad}, nil)
	h := GetProSuggestionsHandler{
		Interactor: mInteractor,
	}
	input := &getProSuggestionsHandlerInput{
		ListID: "1",
	}
	getter := MakeMockInputGetter(input, nil)
	r := h.Execute(getter)

	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: getProSuggestionsHandlerOutput{Ads: []AdsOutput{{ListID: "1", Category: "2020"}}},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestGetProSuggestionsHandlerError(t *testing.T) {
	mInteractor := &mockGetSuggestions{}
	err := fmt.Errorf("err")
	mInteractor.On(
		"GetProSuggestions", mock.Anything, mock.Anything, mock.Anything,
	).Return([]domain.Ad{}, err)

	h := GetProSuggestionsHandler{
		Interactor: mInteractor,
	}
	input := &getProSuggestionsHandlerInput{
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
		"GetProSuggestions", mock.Anything, mock.Anything, mock.Anything,
	).Return([]domain.Ad{}, nil)
	h := GetProSuggestionsHandler{
		Interactor: mInteractor,
	}
	input := &getProSuggestionsHandlerInput{
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
	ads := []domain.Ad{
		{
			ListID:     1,
			CategoryID: 2020,
			Currency:   "uf",
			Price:      1000000,
		},
	}
	h := GetProSuggestionsHandler{
		CurrencySymbol:      "$",
		UnitOfAccountSymbol: "UF",
	}
	r := h.setOutput(ads, []string{})

	expected := getProSuggestionsHandlerOutput{
		Ads: []AdsOutput{
			{
				ListID:   "1",
				Category: "2020",
				Currency: "UF",
				Price:    10000,
			},
		},
	}
	assert.Equal(t, expected, r)
}
func TestGetProSuggestionsHandlerOtherCurrency(t *testing.T) {
	ads := []domain.Ad{
		{
			ListID:     1,
			CategoryID: 2020,
			Currency:   "peso",
			Price:      1000000,
		},
	}
	h := GetProSuggestionsHandler{
		CurrencySymbol:      "$",
		UnitOfAccountSymbol: "UF",
	}
	r := h.setOutput(ads, []string{})

	expected := getProSuggestionsHandlerOutput{
		Ads: []AdsOutput{
			{
				ListID:   "1",
				Category: "2020",
				Currency: "$",
				Price:    1000000,
			},
		},
	}
	assert.Equal(t, expected, r)
}
func TestGetProSuggestionsHandlerOptionalParams(t *testing.T) {
	ads := []domain.Ad{
		{
			ListID:        1,
			CategoryID:    2020,
			Currency:      "peso",
			Price:         1000000,
			PublisherType: domain.Pro,
			AdParams: map[string]string{
				"brand":   "TOYOTA",
				"Mileage": "1234"},
		},
	}
	h := GetProSuggestionsHandler{
		CurrencySymbol:      "$",
		UnitOfAccountSymbol: "UF",
	}
	r := h.setOutput(ads, []string{"publisherType", "mileage", "brand", "regDate", "unknown"})

	expected := getProSuggestionsHandlerOutput{
		Ads: []AdsOutput{
			{
				ListID:        "1",
				Category:      "2020",
				Currency:      "$",
				Price:         1000000,
				PublisherType: "pro",
				Brand:         "TOYOTA",
				Mileage:       "1234",
			},
		},
	}
	assert.Equal(t, expected, r)
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
