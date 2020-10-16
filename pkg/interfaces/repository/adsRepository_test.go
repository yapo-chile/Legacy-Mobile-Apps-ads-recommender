package repository

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/usecases"
)

const getAdTemplateName = "getAd"
const getAdsTemplateName = "getAds"

func TestNewAdsRepository(t *testing.T) {
	mHandler := MockElasticSearchHandler{}
	expected := &adsRepository{
		elasticHandler: &mHandler,
	}
	repo := NewAdsRepository(&mHandler, nil, nil, "", "", 0, 0)
	assert.Equal(t, expected, repo)
	mHandler.AssertExpectations(t)
}

func TestGetAdOK(t *testing.T) {
	mHandler := MockElasticSearchHandler{}
	mDataMapping := MockDataMapping{}
	templateValue, _ := template.New(getAdTemplateName).Parse("test")
	templates := map[string]*template.Template{
		getAdTemplateName: templateValue,
	}
	mDataMapping.On("Get", mock.Anything).Return("test")
	mHandler.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(`
    {"hits" : {"hits" : [{"_source" : {"AdID" : 1,"ListID" : 1, "Subject": "ad testing"}}]}}`, nil)

	repo := adsRepository{
		elasticHandler: &mHandler,
		queryTemplates: templates,
		regionsConf:    &mDataMapping,
	}
	resp, err := repo.GetAd("1")
	expected := domain.Ad{ListID: 1, Subject: "ad testing", URL: "/test/ad_testing_1"}
	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
	mHandler.AssertExpectations(t)
	mDataMapping.AssertExpectations(t)
}

func TestGetAdErr(t *testing.T) {
	mHandler := MockElasticSearchHandler{}
	templateValue, _ := template.New(getAdTemplateName).Parse("test")
	templates := map[string]*template.Template{
		getAdTemplateName: templateValue,
	}
	mHandler.On(
		"Search",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(`{}`, fmt.Errorf(""))

	repo := adsRepository{
		elasticHandler: &mHandler,
		queryTemplates: templates,
	}
	resp, err := repo.GetAd("1")
	expected := domain.Ad{}
	assert.Equal(t, expected, resp)
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
}

func TestGetAdNotEnough(t *testing.T) {
	mHandler := MockElasticSearchHandler{}
	templateValue, _ := template.New(getAdTemplateName).Parse("test")
	templates := map[string]*template.Template{
		getAdTemplateName: templateValue,
	}
	mHandler.On(
		"Search",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(`{}`, nil)

	repo := adsRepository{
		elasticHandler: &mHandler,
		queryTemplates: templates,
	}
	resp, err := repo.GetAd("1")
	expected := domain.Ad{}
	assert.Equal(t, expected, resp)
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
}

func TestGetAdsOK(t *testing.T) {
	mHandler := MockElasticSearchHandler{}
	mDataMapping := MockDataMapping{}
	templateValue, _ := template.New(getAdsTemplateName).Parse("test")
	templates := map[string]*template.Template{
		getAdsTemplateName: templateValue,
	}
	mDataMapping.On("Get", mock.Anything).Return("test")
	mHandler.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(`
    {"hits" : {"hits" : [{"_source" : {"AdID" : 1,"ListID" : 1, "Subject": "ad testing"}}]}}`, nil)

	repo := adsRepository{
		elasticHandler: &mHandler,
		queryTemplates: templates,
		regionsConf:    &mDataMapping,
	}
	options := map[string]string{}
	resp, err := repo.GetAds(options, options, options, options, 1, 0)
	expected := []domain.Ad{{ListID: 1, Subject: "ad testing", URL: "/test/ad_testing_1"}}
	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
	mHandler.AssertExpectations(t)
	mDataMapping.AssertExpectations(t)
}

func TestGetAdsProcessNoTemplate(t *testing.T) {
	templateValue, _ := template.New(getAdsTemplateName).Parse("test")
	templates := map[string]*template.Template{
		getAdsTemplateName: templateValue,
	}
	repo := adsRepository{
		queryTemplates: templates,
	}
	resp, err := repo.getAdsProcess("test2", nil, 0, 0)
	var expected []domain.Ad
	assert.Equal(t, expected, resp)
	assert.Error(t, err)
}

func TestGetAdsProcessUnmarshalErr(t *testing.T) {
	mHandler := MockElasticSearchHandler{}
	templateValue, _ := template.New(getAdsTemplateName).Parse("test")
	templates := map[string]*template.Template{
		getAdsTemplateName: templateValue,
	}
	mHandler.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(`{`, nil)

	repo := adsRepository{
		elasticHandler: &mHandler,
		queryTemplates: templates,
	}
	resp, err := repo.getAdsProcess(getAdsTemplateName, nil, 0, 0)
	var expected []domain.Ad
	assert.Equal(t, expected, resp)
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
}

func TestProcessTemplateErr(t *testing.T) {
	templateValue := template.New(getAdsTemplateName)
	templates := map[string]*template.Template{
		getAdsTemplateName: templateValue,
	}
	repo := adsRepository{
		queryTemplates: templates,
	}
	resp, err := repo.ProcessTemplate(getAdsTemplateName, nil)
	expected := ""
	assert.Equal(t, expected, resp)
	assert.Error(t, err)
}

func TestGetBoolParamOK(t *testing.T) {
	repo := adsRepository{}
	parameters := map[string]string{"Key1": "a", "Key2": "b"}
	resp := repo.getBoolParameters(parameters)
	expected := `{"match": {"Key1": "a"}}, {"match": {"Key2": "b"}}`
	assert.Equal(t, expected, resp)
}

func TestGetFilterParamOK(t *testing.T) {
	repo := adsRepository{}
	parameters := map[string]string{"Key1": "a", "Key2": "b"}
	resp := repo.getFilters(parameters)
	expected := `{"term": {"Key1.keyword": "a"}}, {"term": {"Key2.keyword": "b"}}`
	assert.Equal(t, expected, resp)
}

func TestGetMainImageOK(t *testing.T) {
	repo := adsRepository{imageServerLink: "test/%s/%s/%s"}
	images := []usecases.AdMedia{{ID: 1}}
	resp := repo.getMainImage(images)
	expected := domain.Image{
		Full:   "test/images/00/0000000001",
		Medium: "test/thumbsli/00/0000000001",
		Small:  "test/thumbs/00/0000000001"}
	assert.Equal(t, expected, resp)
}

func TestGetMainImageSeqNoNot0(t *testing.T) {
	repo := adsRepository{imageServerLink: "test/%s/%s/%s"}
	images := []usecases.AdMedia{{ID: 1, SeqNo: 1}}
	resp := repo.getMainImage(images)
	expected := domain.Image{
		Full:   "test/images/00/0000000001",
		Medium: "test/thumbsli/00/0000000001",
		Small:  "test/thumbs/00/0000000001"}
	assert.Equal(t, expected, resp)
}

func TestFillImageOK(t *testing.T) {
	repo := adsRepository{imageServerLink: "test/%s/%s/%s"}
	resp := repo.fillImage(1)
	expected := domain.Image{
		Full:   "test/images/00/0000000001",
		Medium: "test/thumbsli/00/0000000001",
		Small:  "test/thumbs/00/0000000001"}
	assert.Equal(t, expected, resp)
}
