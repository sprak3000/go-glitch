package glitch

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPProblem should be used as the response in case of an error during an HTTP request.
// It implements the https://datatracker.ietf.org/doc/rfc7807 spec with an addition code
// field which is meant to be machine-readable and give clients enough information to handle the error appropriately.
// swagger:model HTTPProblem
type HTTPProblem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Code     string `json:"code,omitempty"`
}

// Error provides a human-readable string outlining the error details
func (h HTTPProblem) Error() string {
	return fmt.Sprintf("HTTPProblem: [%d - %s] - %s - %s", h.Status, h.Code, h.Title, h.Detail)
}

// ResponseToProblem decodes an HTTP response into a problem structure
func ResponseToProblem(res *http.Response) (*HTTPProblem, error) {
	prob := new(HTTPProblem)
	dec := json.NewDecoder(res.Body)
	return prob, dec.Decode(prob)
}

// ValidateProblem validates an HTTP response contains the expected response and status code
func ValidateProblem(res *http.Response, expectedErrorCode string, expectedStatus int) (bool, error) {
	prob, err := ResponseToProblem(res)
	if err != nil {
		return false, err
	}

	return expectedErrorCode == prob.Code && expectedStatus == prob.Status && res.StatusCode == prob.Status, nil
}
