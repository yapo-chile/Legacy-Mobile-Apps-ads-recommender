package usecases

import "github.mpi-internal.com/Yapo/ads-recommender/pkg/domain"

// AdsRepository defines the methods that are available for ad repository
type AdsRepository interface {
	GetAd(listID string) (ad domain.Ad, err error)
	GetAds(
		musts, shoulds, mustsNot, filters, ranges, decay map[string]string,
		queryString []map[string]string,
		size, from int,
	) ([]domain.Ad, error)
}

// AdContactRepo implements ad contact repository functions
type AdContactRepo interface {
	GetAdsPhone(ads []domain.Ad) (phones map[string]string, err error)
}

// IndicatorsRepository defines the methods that a Indicators repository should have
type IndicatorsRepository interface {
	GetUF() (float64, error)
}
