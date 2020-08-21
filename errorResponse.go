package profitwell

import (
	"encoding/json"
	"io"

	"github.com/hatchify/errors"
)

func handleError(r io.Reader) (err error) {
	var resp ErrorResponse
	if err = json.NewDecoder(r).Decode(&resp); err != nil {
		return
	}

	return resp.Error()
}

// ErrorResponse is the response struct for ProfitWell errors
type ErrorResponse struct {
	Detail string `json:"detail"`
}

func (e *ErrorResponse) Error() (err error) {
	return errors.Error(e.Detail)
}
