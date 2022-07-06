package smsman

import _ "embed"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//go:embed countries.json
var countries []byte

const (
	// FivesimAPIEndpoint is the basepoint for 5sim API
	SmsManAPIEndpoint = "http://api.sms-man.com/control/"

	// ANY represents an "any" parameter
	ANY = "any"

	// VERSION is this API wrapper version
	VERSION = "1.0"
)

// Client will perform all the API-related tasks
type Client struct {
	APIKey   string
	Referral string
}

//NewClient get a new Client with a given APIKey
func NewClient(APIKey string) *Client {
	return &Client{APIKey: APIKey}
}

// makeGetRequest performs a simple get request with custom header and query values
func (c *Client) makeGetRequest(url string, queryValues *url.Values) (*http.Response, error) {
	// Creates a client
	client := &http.Client{}
	// Creates a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, err
	}

	queryValues.Set("token", c.APIKey)
	// Encode the query values (if any)
	req.URL.RawQuery = queryValues.Encode()
	return client.Do(req)
}

// BuyActivationNumber performs a "buy activation number" operation by selecting country, operator and product name
// and returns the operation information
func (c *Client) GetNumber(countryId, applicationId string) (*PhoneDetail, error) {

	// Check if any additional query values could be encapsulated
	queryValues := url.Values{}

	queryValues.Add("country_id", countryId)
	queryValues.Add("application_id", applicationId)

	// Make request
	resp, err := c.makeGetRequest(
		fmt.Sprintf("%s/get-number", SmsManAPIEndpoint),
		&queryValues,
	)

	if err != nil {
		return &PhoneDetail{}, err
	}

	// Check status code
	if resp.StatusCode != 200 {
		return &PhoneDetail{}, fmt.Errorf("%s", resp.Status)
	}

	// Read request body
	r, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(r))
	if err != nil {
		resp.Body.Close()
		return &PhoneDetail{}, err
	}
	resp.Body.Close()

	// Unmarshal the body into a struct
	var numberDetail PhoneDetail
	err = json.Unmarshal(r, &numberDetail)
	if err != nil {
		return &PhoneDetail{}, err
	}

	countryCode, phone := c.GetCountry(fmt.Sprintf("+%s", numberDetail.Phone))
	numberDetail.Country = countryCode
	numberDetail.Phone = phone

	return &numberDetail, nil
}

// GetUserInfo returns ID, Email, Balance and rating of the user in a single request
func (c *Client) GetCode(taskid int) (*SmsDetail, error) {

	queryValues := url.Values{}
	queryValues.Add("taskid", strconv.Itoa(taskid))

	// Make request
	resp, err := c.makeGetRequest(
		fmt.Sprintf("%s/get-sms", SmsManAPIEndpoint),
		&queryValues,
	)

	if err != nil {
		return &SmsDetail{}, err
	}

	// Check status code
	if resp.StatusCode != 200 {
		return &SmsDetail{}, fmt.Errorf("%s", resp.Status)
	}

	// Read request body
	r, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(r))
	if err != nil {
		resp.Body.Close()
		return &SmsDetail{}, err
	}
	resp.Body.Close()

	// Unmarshal the body into a struct
	var info SmsDetail
	err = json.Unmarshal(r, &info)
	if err != nil {
		return &SmsDetail{}, err
	}

	return &info, nil
}

func (c *Client) GetCountry(phone string) (conuntryCode, newPhone string) {
	var country map[string]interface{}
	err := json.Unmarshal(countries, &country)
	if err != nil {

	}

	for key, value := range country {
		info := value.(map[string]interface{})
		if strings.Contains(phone, info["code"].(string)) {
			conuntryCode = key
			newPhone = strings.Replace(phone, info["code"].(string), "", -1)
			return
		}
	}

	return "", ""
}
