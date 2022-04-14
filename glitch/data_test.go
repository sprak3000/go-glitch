package glitch

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnit_dataError_Error(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path": {
			validate: func(t *testing.T) {
				e := NewDataError(errors.New("inner err"), "err code", "err msg")
				expect := "Code: [err code] Message: [err msg] Inner error: [inner err]"
				require.Equal(t, expect, e.Error())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}

func TestUnit_dataError_Inner(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path": {
			validate: func(t *testing.T) {
				e := NewDataError(errors.New("inner err"), "err code", "err msg")
				require.Equal(t, errors.New("inner err"), e.Inner())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}

func TestUnit_dataError_Code(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path": {
			validate: func(t *testing.T) {
				e := NewDataError(errors.New("inner err"), "err code", "err msg")
				require.Equal(t, "err code", e.Code())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}

func TestUnit_dataError_Wrap_GetCause(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path": {
			validate: func(t *testing.T) {
				e := NewDataError(errors.New("inner err"), "err code", "err msg").Wrap(NewDataError(nil, "wrapped err code", "wrapped err msg"))
				require.Equal(t, NewDataError(nil, "wrapped err code", "wrapped err msg"), e.GetCause())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}

func TestUnit_FromHTTPProblem(t *testing.T) {
	tests := map[string]struct {
		setupErr    func(t *testing.T) (error, string)
		expectedErr DataError
		validate    func(t *testing.T, expectedErr, actualErr DataError)
	}{
		"base path- inner err not an HTTPProblem": {
			setupErr: func(t *testing.T) (error, string) {
				return NewDataError(errors.New("not an HTTP problem"), "err code", "err msg"), "err msg param"
			},
			expectedErr: NewDataError(NewDataError(errors.New("not an HTTP problem"), "err code", "err msg"), UnknownCode, "err msg param"),
			validate: func(t *testing.T, expectedErr, actualErr DataError) {
				require.Equal(t, expectedErr, actualErr)
			},
		},
		"base path- inner err is an HTTPProblem": {
			setupErr: func(t *testing.T) (error, string) {
				return HTTPProblem{
					Type:     "problem type",
					Title:    "problem title",
					Status:   http.StatusTeapot,
					Detail:   "problem detail",
					Instance: "problem instance",
					Code:     "problem code",
				}, "err msg param"
			},
			expectedErr: NewDataError(HTTPProblem{
				Type:     "problem type",
				Title:    "problem title",
				Status:   http.StatusTeapot,
				Detail:   "problem detail",
				Instance: "problem instance",
				Code:     "problem code",
			}, "problem code", "err msg param"),
			validate: func(t *testing.T, expectedErr, actualErr DataError) {
				require.Equal(t, expectedErr, actualErr)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			e, m := tc.setupErr(t)
			err := FromHTTPProblem(e, m)
			tc.validate(t, tc.expectedErr, err)
		})
	}
}
