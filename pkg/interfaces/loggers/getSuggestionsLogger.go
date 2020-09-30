package loggers

import "github.mpi-internal.com/Yapo/pro-carousel/pkg/usecases"

type getSuggestionsLogger struct {
	logger Logger
}

func (l *getSuggestionsLogger) LimitExceeded(err error) {
	l.logger.Error("%+v", err)
}

func (l *getSuggestionsLogger) MinimumQtyNotEnough(err error) {
	l.logger.Error("%+v", err)
}

func (l *getSuggestionsLogger) ErrorGettingAd(listID string, err error) {
	l.logger.Error("cannot get ad with listID %s with error: %+v", listID, err)
}

func (l *getSuggestionsLogger) ErrorGettingAds(musts, shoulds, mustsNot map[string]string, err error) {
	l.logger.Error("cannot get ads using params %v - %v - %v with err %+v", musts, shoulds, mustsNot, err)
}

func (l *getSuggestionsLogger) NotEnoughAds(listID string, lenAds int) {
	l.logger.Info("cannot get enough ads using listID, just got %d ads", listID, lenAds)
}

// MakeGetSuggestionsLogger sets up a GetSuggestionsLogger instrumented
// via the provided logger
func MakeGetSuggestionsLogger(logger Logger) usecases.GetSuggestionsLogger {
	return &getSuggestionsLogger{
		logger: logger,
	}
}
