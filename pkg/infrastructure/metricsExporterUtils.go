package infrastructure

import (
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	// eventTypeFuncCallDeepness deepness of function invokation for eventType
	eventTypeFuncCallDeepness int = 4
	// eventNameFuncCallDeepness deepness of function invokation for eventName
	eventNameFuncCallDeepness int = 5
	// entityNameFuncCallDeepness deepness of function invokation for entityName
	entityNameFuncCallDeepness int = 6
	// stackTraceInitialBufferCapacity initial buffer capacity for stack trace
	stackTraceInitialBufferCapacity int = 20
)

// getEventType gets the event type to export to prometheus
// example: fiboRepo.logRepositoryError calls to logger.Error("...")
// getEventType gets "error" as event type
func getEventType() string {
	return toSnakeCase(getFuncName(eventTypeFuncCallDeepness))
}

// getEntityName gets the entity name to export to prometheus.
// example: fiboRepo.logRepositoryError calls to logger.Error("...")
// getEventType gets "fibo_repo" as entity name
func getEntityName() string {
	return toSnakeCase(getEntity(entityNameFuncCallDeepness))
}

var loggerReplacer = regexp.MustCompile(`(_?log(ger)?_?)`) // nolint: gochecknoglobals

// getEventName gets the event name to export to prometheus.
// example: fiboRepo.logRepositoryError calls to logger.Error("...")
// getEventType gets "repository_error" as event name removing "log" or "logger"
func getEventName() string {
	loggerName := toSnakeCase(getFuncName(eventNameFuncCallDeepness))
	return loggerReplacer.ReplaceAllString(loggerName, "")
}

func getFuncName(deepness int) string {
	nameFull := funcName(deepness, stackTraceInitialBufferCapacity)
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

var entityRgx = regexp.MustCompile(`\(\*\w+\)`) // nolint: gochecknoglobals
var wordsOnly = regexp.MustCompile(`\w+`)       // nolint: gochecknoglobals

func getEntity(deepness int) string {
	nameFull := funcName(deepness, stackTraceInitialBufferCapacity)
	entityName := wordsOnly.FindString(entityRgx.FindString(nameFull))
	if entityName == "" {
		return wordsOnly.FindString(filepath.Ext(nameFull))
	}
	return entityName
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)") // nolint: gochecknoglobals

func toSnakeCase(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}

// funcName returns the last function name of invocations on the calling goroutine's stack.
// The stack trace can be skipped using deepness parameter. BufferInitialCapacity gives
// the initial capacity to record in trace stack slice.
func funcName(deepness, bufferInitialCapacity int) string {
	pc := make([]uintptr, bufferInitialCapacity)
	n := runtime.Callers(deepness, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}
