package usecases

import (
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/domain"
)

// GetSuggestionsInteractor defines the available methods for this interactor
type GetSuggestionsInteractor interface {
	// GetProSuggestions will get all suggestions for the given listID
	GetProSuggestions(
		listID string,
		optionalParams []string,
		size, from int,
		carousel string,
	) (ads []domain.Ad, err error)
}
