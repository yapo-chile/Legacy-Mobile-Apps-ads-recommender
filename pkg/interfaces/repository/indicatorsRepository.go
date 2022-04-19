package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.mpi-internal.com/Yapo/ads-recommender/pkg/usecases"
)

// indicatorsRepository loan settings datasource
type indicatorsRepository struct {
	HTTPCachedHandler HTTPCachedHandler
	UFPath            string
	DefaultValue      float64
}

// NewIndicatorsRepository returns a indicatorsRepository instance
func NewIndicatorsRepository(
	httpCachedHandler HTTPCachedHandler,
	ufPath string,
	defaultValue int,
) usecases.IndicatorsRepository {
	return &indicatorsRepository{
		HTTPCachedHandler: httpCachedHandler,
		UFPath:            ufPath,
		DefaultValue:      float64(defaultValue),
	}
}

// GetUF get UF value
func (repo *indicatorsRepository) GetUF() (float64, error) {
	t := time.Now()
	dateStr := fmt.Sprintf("%02d-%02d-%d", t.Day(), t.Month(), t.Year())
	request := repo.HTTPCachedHandler.NewRequest().
		SetMethod("GET").
		SetPath(repo.UFPath + dateStr)
	response, err := repo.HTTPCachedHandler.Send(request)
	if err == nil && response != nil {
		var ufAPIResponse usecases.UFApiResponse
		b := []byte(response.(string))
		err = json.Unmarshal(b, &ufAPIResponse)
		if err == nil {
			if len(ufAPIResponse.Sets) > 0 {
				return ufAPIResponse.Sets[0].Value, nil
			}
			return repo.DefaultValue, fmt.Errorf(usecases.ErrGetUF)
		}
	}
	return repo.DefaultValue, err
}
