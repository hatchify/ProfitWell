package profitwell

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const (
	v1Host = "https://www2.profitwell.com"
	v1API  = "/dotjs/v1"

	v2Host = "https://api.profitwell.com"
	v2API  = "/v2"
)

const (
	// SetCustomerActionEndpoint is the endpoint for setting a customer action
	SetCustomerActionEndpoint = "/quests/customer"
)

// New will return a new instance of ProfitWell
func New(authToken string) (pp *ProfitWell, err error) {
	var p ProfitWell
	if p.v1URL, err = url.Parse(v1Host); err != nil {
		err = fmt.Errorf("error parsing v1 host: %v", err)
		return
	}

	if p.v2URL, err = url.Parse(v2Host); err != nil {
		err = fmt.Errorf("error parsing v2 host: %v", err)
		return
	}

	p.authToken = authToken
	pp = &p
	return
}

// ProfitWell is a manager for the ProfitWell SDK
type ProfitWell struct {
	hc http.Client

	v1URL *url.URL
	v2URL *url.URL

	// ProfitWell authentication token
	// Note: This determines which account is performing actions within their API
	authToken string
}

func (p *ProfitWell) getV1URL(endpoint string) (u url.URL) {
	u = *p.v1URL
	u.Path = path.Join(v1API, endpoint)
	return
}

func (p *ProfitWell) getV2URL(endpoint string) (u url.URL) {
	u = *p.v2URL
	u.Path = path.Join(v2API, endpoint)
	return
}

func (p *ProfitWell) request(method, endpoint string, response interface{}) (err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, endpoint, nil); err != nil {
		return
	}

	// Set authorization token within header
	// Note: The same authorization is used for both v1 and v2
	req.Header.Set("Authorization", p.authToken)

	var resp *http.Response
	if resp, err = p.hc.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return handleError(resp.Body)
	}

	if response == nil {
		// No response was expected, return before parsing response body
		return
	}

	// Parse response body as JSON
	return json.NewDecoder(resp.Body).Decode(response)
}

// SetUserAction will set a recent action for a given user email
func (p *ProfitWell) SetUserAction(userEmail string) (err error) {
	u := p.getV1URL(SetCustomerActionEndpoint)
	q := url.Values{}
	q.Set("user_email", userEmail)
	u.RawQuery = q.Encode()

	return p.request("GET", u.String(), nil)
}
