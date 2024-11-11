package helpers

import "github.com/rs/zerolog/log"

func CoreLogError(method string, message string, isFatal bool) {
	if isFatal {
		log.Fatal().
			Str("method", method).
			Msg(message)
	} else {
		log.Error().
			Str("method", method).
			Msg(message)
	}
}

func CoreLogInfo(method string, message string) {
	log.Info().
		Str("method", method).
		Msg(message)
}

func HttpLogNotFound(method string, message string) {
	httpLogError(method, "Not found", message)
}

func HttpLogBadRequest(method string, message string) {
	httpLogError(method, "Bad request", message)
}

func HttpLogISR(method string, message string) {
	httpLogError(method, "Internal server error", message)
}

func HttpLogConflict(method string, message string) {
	httpLogError(method, "Conflict", message)
}

func httpLogError(method string, errorType string, message string) {
	log.Error().
		Str("method", method).
		Str("errorType", errorType).
		Msg(message)
}

func HttpLogInfo(method string, message string) {
	log.Error().
		Str("method", method).
		Msg(message)
}
