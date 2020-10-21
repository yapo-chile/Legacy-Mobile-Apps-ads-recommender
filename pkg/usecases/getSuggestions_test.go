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

type mockAdsRepository struct {
	mock.Mock
}

func (m *mockAdsRepository) GetAd(listID string) (domain.Ad, error) {
	args := m.Called(listID)
	return args.Get(0).(domain.Ad), args.Error(1)
}
func (m *mockAdsRepository) GetAds(
	musts, shoulds, mustsNot, filters map[string]string, size, from int,
) ([]domain.Ad, error) {
	args := m.Called(musts, shoulds, mustsNot, filters)
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

func TestGetProSuggestionsOK(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	ads := []domain.Ad{{ListID: 2, Category: "test"}}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ads, nil)
	i := GetSuggestions{
		SuggestionsRepo: &mAdsRepo,
		MinDisplayedAds: 1,
		MaxDisplayedAds: 1,
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{}, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
}

func TestGetProSuggestionsMaxDisplayedAds(t *testing.T) {
	mAdsRepo := mockAdsRepository{}
	mLogger := mockGetSuggestionsLogger{}
	ad := domain.Ad{ListID: 1, Category: "test"}
	ads := []domain.Ad{{ListID: 2, Category: "test"}}
	mAdsRepo.On("GetAd", mock.Anything).Return(ad, nil)
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ads, nil)
	mLogger.On("LimitExceeded", mock.Anything, mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo: &mAdsRepo,
		MinDisplayedAds: 1,
		MaxDisplayedAds: 1,
		RequestedAdsQty: 2,
		Logger:          &mLogger,
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{}, 2, 0)
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
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ads, nil)
	mLogger.On("MinimumQtyNotEnough", mock.Anything, mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo: &mAdsRepo,
		MinDisplayedAds: 2,
		MaxDisplayedAds: 2,
		RequestedAdsQty: 2,
		Logger:          &mLogger,
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{}, 1, 0)
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
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ads, nil)
	mLogger.On("NotEnoughAds", mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo: &mAdsRepo,
		MinDisplayedAds: 2,
		MaxDisplayedAds: 2,
		RequestedAdsQty: 2,
		Logger:          &mLogger,
	}
	expected := []domain.Ad{}
	output, err := i.GetProSuggestions("1", []string{}, 2, 0)
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
		SuggestionsRepo: &mAdsRepo,
		MinDisplayedAds: 1,
		MaxDisplayedAds: 2,
		Logger:          &mLogger,
	}
	var expected []domain.Ad
	output, err := i.GetProSuggestions("1", []string{}, 1, 0)
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
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]domain.Ad{}, fmt.Errorf(""))
	mLogger.On("ErrorGettingAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo: &mAdsRepo,
		MinDisplayedAds: 1,
		MaxDisplayedAds: 1,
		RequestedAdsQty: 2,
		Logger:          &mLogger,
	}
	expected := []domain.Ad{}
	output, err := i.GetProSuggestions("1", []string{}, 1, 0)
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
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ads, nil)
	mAdContactRepo.On("GetAdsPhone", mock.Anything).Return(phones, nil)
	i := GetSuggestions{
		SuggestionsRepo: &mAdsRepo,
		AdContact:       &mAdContactRepo,
		MinDisplayedAds: 1,
		MaxDisplayedAds: 2,
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{"phonelink"}, 1, 0)
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
	mAdsRepo.On("GetAds", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ads, nil)
	mAdContactRepo.On("GetAdsPhone", mock.Anything).Return(phones, fmt.Errorf("error"))
	mLogger.On("ErrorGettingAdsContact", mock.Anything, mock.Anything)
	i := GetSuggestions{
		SuggestionsRepo: &mAdsRepo,
		AdContact:       &mAdContactRepo,
		MinDisplayedAds: 1,
		MaxDisplayedAds: 2,
		Logger:          &mLogger,
	}
	expected := ads
	output, err := i.GetProSuggestions("1", []string{"phonelink"}, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	mAdsRepo.AssertExpectations(t)
	mAdContactRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetShouldsParamsOK(t *testing.T) {
	ad := domain.Ad{AdParams: map[string]string{"a": "test", "b": "test2"}}
	suggestionsParams := []string{"a", "c", "b"}
	expected := map[string]string{"Params.a": "test", "Params.b": "test2"}
	result := getShouldsParams(ad, suggestionsParams)
	assert.Equal(t, expected, result)
}
