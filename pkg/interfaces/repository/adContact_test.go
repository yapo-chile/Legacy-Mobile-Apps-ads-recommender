package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/domain"
)

func TestNewAdContactRepository(t *testing.T) {
	mHandler := mockHTTPHandler{}
	expected := &AdContactRepository{
		handler: &mHandler,
	}
	repo := NewAdContactRepository(&mHandler, "")
	assert.Equal(t, expected, repo)
	mHandler.AssertExpectations(t)
}

func TestGetAdsPhoneOK(t *testing.T) {
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}

	mRequest.On("SetPath", mock.AnythingOfType("string")).Return(&mRequest)
	mRequest.On("SetMethod", "GET").Return(&mRequest)
	mRequest.On("SetBody").Return(&mRequest)
	mHandler.On("NewRequest").Return(&mRequest)
	mHandler.On("Send", &mRequest).Return(
		`{"1": "www.test1.cl"}`, nil)

	expected := map[string]string{"1": "www.test1.cl"}

	repo := AdContactRepository{
		handler: &mHandler,
	}

	resp, err := repo.GetAdsPhone([]domain.Ad{{ListID: 1}, {ListID: 2}})
	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetAdsPhoneReqErr(t *testing.T) {
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}

	mRequest.On("SetPath", mock.AnythingOfType("string")).Return(&mRequest)
	mRequest.On("SetMethod", "GET").Return(&mRequest)
	mRequest.On("SetBody").Return(&mRequest)

	mHandler.On("NewRequest").Return(&mRequest)
	mHandler.On("Send", &mRequest).Return(
		`{"1": "www.test1.cl"}`, fmt.Errorf("error"))

	repo := AdContactRepository{
		handler: &mHandler,
	}

	_, err := repo.GetAdsPhone([]domain.Ad{{ListID: 1}, {ListID: 2}})
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetAdsPhoneUnMarshalErr(t *testing.T) {
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}

	mRequest.On("SetPath", mock.AnythingOfType("string")).Return(&mRequest)
	mRequest.On("SetMethod", "GET").Return(&mRequest)
	mRequest.On("SetBody").Return(&mRequest)
	mHandler.On("NewRequest").Return(&mRequest)
	mHandler.On("Send", &mRequest).Return(`{`, nil)

	repo := AdContactRepository{
		handler: &mHandler,
	}

	_, err := repo.GetAdsPhone([]domain.Ad{{ListID: 1}, {ListID: 2}})
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetAdsPhoneNoAdsErr(t *testing.T) {
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}

	mRequest.On("SetPath", mock.AnythingOfType("string")).Return(&mRequest)
	mRequest.On("SetMethod", "GET").Return(&mRequest)
	mRequest.On("SetBody").Return(&mRequest)
	mHandler.On("NewRequest").Return(&mRequest)
	mHandler.On("Send", &mRequest).Return(`{}`, nil)

	repo := AdContactRepository{
		handler: &mHandler,
	}

	_, err := repo.GetAdsPhone([]domain.Ad{{ListID: 1}, {ListID: 2}})
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}
