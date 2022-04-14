package glitch

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnit_HTTPProblem_Error(t *testing.T) {
	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path": {
			validate: func(t *testing.T) {
				h := HTTPProblem{
					Type:     "problem type",
					Title:    "problem title",
					Status:   http.StatusTeapot,
					Detail:   "problem detail",
					Instance: "problem instance",
					Code:     "problem code",
				}

				expect := "HTTPProblem: [418 - problem code] - problem title - problem detail"
				require.Equal(t, expect, h.Error())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}

func TestUnit_ResponseToProblem(t *testing.T) {
	tests := map[string]struct {
		setupResp       func(t *testing.T) *http.Response
		expectedProblem HTTPProblem
		validate        func(t *testing.T, expectedProblem, actualProblem *HTTPProblem, actualErr error)
	}{
		"base path": {
			setupResp: func(t *testing.T) *http.Response {
				b := bytes.NewReader([]byte(`{"type":"problem type","title":"problem title","status":418,"detail":"problem detail","instance":"problem instance","code":"problem code"}`))
				return &http.Response{
					Body: io.NopCloser(b),
				}
			},
			expectedProblem: HTTPProblem{
				Type:     "problem type",
				Title:    "problem title",
				Status:   http.StatusTeapot,
				Detail:   "problem detail",
				Instance: "problem instance",
				Code:     "problem code",
			},
			validate: func(t *testing.T, expectedProblem, actualProblem *HTTPProblem, actualErr error) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedProblem, actualProblem)
			},
		},
		"exceptional path- decode err": {
			setupResp: func(t *testing.T) *http.Response {
				b := bytes.NewReader([]byte(`{malformed "json"`))
				return &http.Response{
					Body: io.NopCloser(b),
				}
			},
			validate: func(t *testing.T, expectedProblem, actualProblem *HTTPProblem, actualErr error) {
				require.Error(t, actualErr)
				require.Equal(t, "invalid character 'm' looking for beginning of object key string", actualErr.Error())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := tc.setupResp(t)
			p, err := ResponseToProblem(r)
			tc.validate(t, &tc.expectedProblem, p, err)
		})
	}
}

func TestUnit_ValidateProblem(t *testing.T) {
	tests := map[string]struct {
		setupResp       func(t *testing.T) *http.Response
		expectedErrCode string
		expectedStatus  int
		expectedResult  bool
		validate        func(t *testing.T, expectedResult, actualResult bool, actualErr error)
	}{
		"base path- response contains expected problem data": {
			setupResp: func(t *testing.T) *http.Response {
				b := bytes.NewReader([]byte(`{"type":"problem type","title":"problem title","status":418,"detail":"problem detail","instance":"problem instance","code":"problem code"}`))
				return &http.Response{
					StatusCode: http.StatusTeapot,
					Body:       io.NopCloser(b),
				}
			},
			expectedErrCode: "problem code",
			expectedStatus:  http.StatusTeapot,
			expectedResult:  true,
			validate: func(t *testing.T, expectedResult, actualResult bool, actualErr error) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedResult, actualResult)
			},
		},
		"base path- response does not contain expected problem error code in body": {
			setupResp: func(t *testing.T) *http.Response {
				b := bytes.NewReader([]byte(`{"type":"problem type","title":"problem title","status":418,"detail":"problem detail","instance":"problem instance","code":"move along"}`))
				return &http.Response{
					StatusCode: http.StatusTeapot,
					Body:       io.NopCloser(b),
				}
			},
			expectedErrCode: "problem code",
			expectedStatus:  http.StatusTeapot,
			expectedResult:  false,
			validate: func(t *testing.T, expectedResult, actualResult bool, actualErr error) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedResult, actualResult)
			},
		},
		"base path- response does not contain expected problem status code in body": {
			setupResp: func(t *testing.T) *http.Response {
				b := bytes.NewReader([]byte(`{"type":"problem type","title":"problem title","status":418,"detail":"problem detail","instance":"problem instance","code":"problem code"}`))
				return &http.Response{
					StatusCode: http.StatusTeapot,
					Body:       io.NopCloser(b),
				}
			},
			expectedErrCode: "problem code",
			expectedStatus:  http.StatusBadRequest,
			expectedResult:  false,
			validate: func(t *testing.T, expectedResult, actualResult bool, actualErr error) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedResult, actualResult)
			},
		},
		"base path- response does not contain expected problem status code in its header": {
			setupResp: func(t *testing.T) *http.Response {
				b := bytes.NewReader([]byte(`{"type":"problem type","title":"problem title","status":418,"detail":"problem detail","instance":"problem instance","code":"problem code"}`))
				return &http.Response{
					StatusCode: http.StatusTeapot,
					Body:       io.NopCloser(b),
				}
			},
			expectedErrCode: "problem code",
			expectedStatus:  http.StatusBadRequest,
			expectedResult:  false,
			validate: func(t *testing.T, expectedResult, actualResult bool, actualErr error) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedResult, actualResult)
			},
		},
		"exceptional path- decode err": {
			setupResp: func(t *testing.T) *http.Response {
				b := bytes.NewReader([]byte(`{malformed "json"`))
				return &http.Response{
					StatusCode: http.StatusTeapot,
					Body:       io.NopCloser(b),
				}
			},
			validate: func(t *testing.T, expectedResult, actualResult bool, actualErr error) {
				require.Error(t, actualErr)
				require.Equal(t, "invalid character 'm' looking for beginning of object key string", actualErr.Error())
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := tc.setupResp(t)
			res, err := ValidateProblem(r, tc.expectedErrCode, tc.expectedStatus)
			tc.validate(t, tc.expectedResult, res, err)
		})
	}
}
