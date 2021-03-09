package loggers

import "github.mpi-internal.com/Yapo/ads-recommender/pkg/usecases"

type getSuggestionsLogger struct {
	logger Logger
}

// LimitExceeded logs ads limit exceeded
func (l *getSuggestionsLogger) LimitExceeded(size, maxDisplayedAds, defaultAdsQty int) {
	l.logger.Info(
		"requesting %d ads but the limit to display is %d, setting size on %d",
		size, maxDisplayedAds, defaultAdsQty,
	)
}

// MinimumQtyNotEnough logs when the minimum ads are not enough
func (l *getSuggestionsLogger) MinimumQtyNotEnough(size, minDisplayedAds, defaultAdsQty int) {
	l.logger.Info(
		"requesting %d ads but the minimum ads quantity to display is %d, setting size on %d",
		size, minDisplayedAds, defaultAdsQty,
	)
}

// ErrorGettingAd logs when cannot get ad
func (l *getSuggestionsLogger) ErrorGettingAd(listID string, err error) {
	l.logger.Error("cannot get ad with listID %s with error: %+v", listID, err)
}

// ErrorGettingAds logs when cannot get ads
func (l *getSuggestionsLogger) ErrorGettingAds(musts, shoulds, mustsNot map[string]string, err error) {
	l.logger.Error("cannot get ads using params %v - %v - %v with err %+v", musts, shoulds, mustsNot, err)
}

// ErrorGettingUF logs when cannot get uf value
func (l *getSuggestionsLogger) ErrorGettingUF(err error) {
	l.logger.Error("cannot get uf value: %+v", err)
}

// NotEnoughAds logs when ads returned are not enough
func (l *getSuggestionsLogger) NotEnoughAds(listID string, lenAds int) {
	l.logger.Info("cannot get enough ads using listID %s, just got %d ads", listID, lenAds)
}

// ErrorGettingAdsContact logs when cannot get ads contact
func (l *getSuggestionsLogger) ErrorGettingAdsContact(listID string, err error) {
	l.logger.Error("cannot get ads contact with listID %s with error: %+v", listID, err)
}

// InvalidCarousel logs when carousel is not valid
func (l *getSuggestionsLogger) InvalidCarousel(carousel string) {
	l.logger.Warn("carousel '%s' not found", carousel)
}

// MakeGetSuggestionsLogger sets up a GetSuggestionsLogger instrumented
// via the provided logger
func MakeGetSuggestionsLogger(logger Logger) usecases.GetSuggestionsLogger {
	return &getSuggestionsLogger{
		logger: logger,
	}
}
