package debugging

import "time"

type HistoricalError struct {
	error
	time.Time
}

var (
	errStore = make([]HistoricalError, 0)
	errInput = make(chan HistoricalError)
)

func init() {
	go func() {
		for err := range errInput {
			errStore = append(errStore, err)
		}
	}()
}

func RegisterError(err error) {
	if err == nil {
		return
	}
	go func() {
		errInput <- HistoricalError{err, time.Now()}
	}()
}

func Errors() []HistoricalError {
	result := make([]HistoricalError, len(errStore))
	copy(result, errStore)
	return result
}
