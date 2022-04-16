package postgres

import (
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/sprak3000/go-glitch/glitch"
)

func TestUnit_ToDataError(t *testing.T) {
	tests := map[string]struct {
		innerErr    error
		msg         string
		expectedErr glitch.DataError
		validate    func(t *testing.T, expectedErr, actualErr glitch.DataError)
	}{
		"base path- nil inner error": {
			validate: func(t *testing.T, expectedErr, actualErr glitch.DataError) {
				require.Equal(t, expectedErr, actualErr)
			},
		},
		"base path- inner error not a pq.Error compatible type": {
			innerErr:    errors.New("err"),
			msg:         "err msg",
			expectedErr: glitch.NewDataError(errors.New("err"), glitch.UnknownCode, "err msg"),
			validate: func(t *testing.T, expectedErr, actualErr glitch.DataError) {
				require.Equal(t, expectedErr, actualErr)
			},
		},
		"base path- inner error a pq.Error compatible type": {
			innerErr:    &pq.Error{Code: "pq err code"},
			msg:         "err msg",
			expectedErr: glitch.NewDataError(&pq.Error{Code: "pq err code"}, "pq err code", "err msg"),
			validate: func(t *testing.T, expectedErr, actualErr glitch.DataError) {
				require.Equal(t, expectedErr, actualErr)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ToDataError(tc.innerErr, tc.msg)
			tc.validate(t, tc.expectedErr, err)
		})
	}
}
