package usecases

import (
	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/domain"
)

// GetSuggestionsInteractor defines the available methods for this interactor
type GetSuggestionsInteractor interface {
	// GetSuggestions will get all suggestions for the given listID
	GetSuggestions(
		listID string,
		optionalParams []string,
		size, from int,
		carouselType string,
	) (ads []domain.Ad, err error)
}
