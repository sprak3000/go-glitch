package postgres

import (
	"github.com/healthimation/go-glitch/glitch"
	"github.com/lib/pq"
)

// FromPQ will convert a lib/pq error into a DataError
func FromPQ(inner error, msg string) glitch.DataError {
	if inner == nil {
		return nil
	}
	if pqErr, ok := inner.(*pq.Error); ok {
		return glitch.NewDataError(inner, string(pqErr.Code), msg)
	}
	return glitch.NewDataError(inner, glitch.UnknownCode, msg)
}
