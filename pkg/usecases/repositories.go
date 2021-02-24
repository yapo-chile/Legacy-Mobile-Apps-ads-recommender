package usecases

import "github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"

// AdsRepository defines the methods that are available for ad repository
type AdsRepository interface {
	GetAd(listID string) (ad domain.Ad, err error)
	GetAds(
		musts, shoulds, mustsNot, filters, decay map[string]string,
		ranges map[string]string,
		size, from int,
	) ([]domain.Ad, error)
}

// AdContactRepo implements ad contact repository functions
type AdContactRepo interface {
	GetAdsPhone(ads []domain.Ad) (phones map[string]string, err error)
}
