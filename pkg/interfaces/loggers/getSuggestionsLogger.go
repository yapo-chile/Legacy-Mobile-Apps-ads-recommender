package loggers

import "github.mpi-internal.com/Yapo/pro-carousel/pkg/usecases"

type getSuggestionsLogger struct {
	logger Logger
}

func (l *getSuggestionsLogger) LimitExceeded(size, maxDisplayedAds, defaultAdsQty int) {
	l.logger.Info(
		"requesting %d ads but the limit to display is %d, setting size on %d",
		size, maxDisplayedAds, defaultAdsQty,
	)
}

func (l *getSuggestionsLogger) MinimumQtyNotEnough(size, minDisplayedAds, defaultAdsQty int) {
	l.logger.Info(
		"requesting %d ads but the minimum ads quantity to display is %d, setting size on %d",
		size, minDisplayedAds, defaultAdsQty,
	)
}

func (l *getSuggestionsLogger) ErrorGettingAd(listID string, err error) {
	l.logger.Error("cannot get ad with listID %s with error: %+v", listID, err)
}

func (l *getSuggestionsLogger) ErrorGettingAds(musts, shoulds, mustsNot map[string]string, err error) {
	l.logger.Error("cannot get ads using params %v - %v - %v with err %+v", musts, shoulds, mustsNot, err)
}

func (l *getSuggestionsLogger) NotEnoughAds(listID string, lenAds int) {
	l.logger.Info("cannot get enough ads using listID %s, just got %d ads", listID, lenAds)
}

func (l *getSuggestionsLogger) ErrorGettingAdsContact(listID string, err error) {
	l.logger.Error("cannot get ads contact with listID %s with error: %+v", listID, err)
}

func (l *getSuggestionsLogger) InvalidCarousel(carousel string) {
	l.logger.Error("carousel '%s' not found", carousel)
}

// MakeGetSuggestionsLogger sets up a GetSuggestionsLogger instrumented
// via the provided logger
func MakeGetSuggestionsLogger(logger Logger) usecases.GetSuggestionsLogger {
	return &getSuggestionsLogger{
		logger: logger,
	}
}
