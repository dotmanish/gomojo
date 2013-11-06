// Copyright 2013 Manish Malik (manishmalik.name)
// All rights reserved.
// Use of this source code is governed by a BSD (3-Clause) License
// that can be found in the LICENSE file.
//
// I initially wrote bits of this for Instamojo integration with an app backend.
// Now it's a Golang package that can be used to interact with the API.
// The latest is at https://github.com/dotmanish/gomojo
//
// Check gomojo-tool.go for a comprehensive sample of API usage.
//
// Currently available APIs:
//
// Initialization:
// 		InitGomojoWithAuthToken
// 		InitGomojoWithUserPass
//
// Main APIs:
// 		ListOffers
//		GetOfferDetails
// 		GetNewAuthToken
// 		DeleteAuthToken
//
// Helper Functions:
// 		GetCurrentAuthToken
//

package gomojo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ListOffersResponse: represents response of 'offer' API
type ListOffersResponse struct {
	Offers  []Offer `json:"offers"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
}

// OfferDetailsResponse: represents response of 'offer' details API
type OfferDetailsResponse struct {
	Offer   Offer  `json:"offer"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AuthResponse: represents response of 'auth' POST API
type AuthResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// DeAuthResponse: represents response of 'auth' DELETE API
type DeAuthResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Offer: represents one Offer object (in list of offers)
// This is an amalgamation of the fields received as a
// result of various APIs (Offers List / Offer Details)
// Note that Offers List API doesn't populate everything.
type Offer struct {
	ShortURL       string `json:"shorturl"`
	Title          string `json:"title"`
	Slug           string `json:"slug"`
	Status         string `json:"status"`
	Description    string `json:"description"`
	Currency       string `json:"currency"`
	BasePrice      string `json:"base_price"`
	Quantity       string `json:"quantity"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	Timezone       string `json:"timezone"`
	Venue          string `json:"venue"`
	RedirectURL    string `json:"redirect_url"`
	Note           string `json:"note"`
	FileUploadJSON string `json:"file_upload_json"`
	CoverImageJSON string `json:"cover_image_json"`
}

var gomojo_app_id, gomojo_auth_token, gomojo_api_ver, gomojo_version string
var gomojo_username, gomojo_password string
var gomojo_init_done bool

// InitGomojoWithAuthToken: Initialize gomojo with Auth Token
// Inputs: (API version string, App ID string, Auth Token string)
func InitGomojoWithAuthToken(api_ver, app_id, auth_token string) {

	if api_ver != "" && app_id != "" && auth_token != "" {
		gomojo_app_id = app_id
		gomojo_auth_token = auth_token
		gomojo_api_ver = api_ver

		gomojo_init_done = true
	}
}

// initGomojoWithUserPass: Initialize gomojo with username/password
// Inputs: (API version string, App ID string, Username string, Password string)
// Note: The username/password are internally stored in gomojo variables
// until the first attempt to call an API. Afterwards, these internal variables
// are blanked out.
func InitGomojoWithUserPass(api_ver, app_id, username, password string) {

	if api_ver != "" && app_id != "" && username != "" && password != "" {
		gomojo_app_id = app_id
		gomojo_api_ver = api_ver
		// These variables are reset to blank once they're used
		gomojo_username = username
		gomojo_password = password

		gomojo_init_done = true
	}
}

// GetCurrentAuthToken: returns the current Auth Token
// Inputs: None
// Returns (Auth Token string)
func GetCurrentAuthToken() string {
	return gomojo_auth_token
}

// SetCurrentAuthToken: sets the current Auth Token
// Inputs: (Auth Token string)
func SetCurrentAuthToken(auth_token string) {
	gomojo_auth_token = auth_token
}

// ListOffers: retrieves the list of all offers created under the given App(ID)
// Inputs: None
// Returns: (Offer object array, API success bool, Message string)
func ListOffers() ([]Offer, bool, string) {

	jsonobj := new(ListOffersResponse)

	if gomojo_init_done {

		api_result := callAPI("listoffers", "")

		jsonerr := json.Unmarshal([]byte(api_result), jsonobj)

		if jsonerr != nil {
			jsonobj.Message = "Invalid JSON: " + jsonerr.Error()
		}
	} else {
		jsonobj.Message = "Please call gomojo.InitGomojoWithAuthToken() or gomojo.InitGomojoWithUserPass() first."
	}

	return jsonobj.Offers, jsonobj.Success, jsonobj.Message
}

// GetOfferDetails: retrieves the details of a particular offer
// Inputs: (Offer-Slug string)
// Returns: (Offer object, API success bool, Message string)
func GetOfferDetails(offer_slug string) (Offer, bool, string) {

	jsonobj := new(OfferDetailsResponse)

	if gomojo_init_done {

		api_result := callAPI("offerdetails", offer_slug)

		jsonerr := json.Unmarshal([]byte(api_result), jsonobj)

		if jsonerr != nil {
			jsonobj.Message = "Invalid JSON: " + jsonerr.Error()
		}
	} else {
		jsonobj.Message = "Please call gomojo.InitGomojoWithAuthToken() or gomojo.InitGomojoWithUserPass() first."
	}

	return jsonobj.Offer, jsonobj.Success, jsonobj.Message
}

// GetNewAuthToken: gets a new Auth Token
// Inputs: (Username string, Password string)
// Returns: (Auth Token string, API success bool, Message string)
func GetNewAuthToken(username, password string) (string, bool, string) {

	jsonobj := new(AuthResponse)

	if gomojo_init_done {

		api_result := callAPI("auth", "username="+username+"&password="+password)

		jsonerr := json.Unmarshal([]byte(api_result), jsonobj)

		if jsonerr != nil {
			jsonobj.Message = "Invalid JSON: " + jsonerr.Error()
		}
	} else {
		jsonobj.Message = "Please call gomojo.InitGomojoWithAuthToken() or gomojo.InitGomojoWithUserPass() first."
	}

	return jsonobj.Token, jsonobj.Success, jsonobj.Message
}

// DeleteAuthToken: deletes an existing Auth Token
// Inputs: (Auth Token string)
// Returns: (API success bool, Message string)
func DeleteAuthToken(auth_token string) (bool, string) {

	jsonobj := new(DeAuthResponse)

	if gomojo_init_done {

		api_result := callAPI("deauth", auth_token)

		jsonerr := json.Unmarshal([]byte(api_result), jsonobj)

		if jsonerr != nil {
			jsonobj.Message = "Invalid JSON: " + jsonerr.Error()
		}
	} else {
		jsonobj.Message = "Please call gomojo.InitGomojoWithAuthToken() or gomojo.InitGomojoWithUserPass() first."
	}

	return jsonobj.Success, jsonobj.Message
}

// callAPI: Internal function handling the REST API
// Valid apicall values: "auth", "deauth", "listoffers"
func callAPI(apicall, apidata string) string {

	// Check if we have auth token available.
	// If not, let's first authenticate and retrieve it.
	if apicall != "auth" && gomojo_auth_token == "" {

		new_auth_token, new_auth_success, _ := GetNewAuthToken(gomojo_username, gomojo_password)

		gomojo_username = ""
		gomojo_password = ""

		if new_auth_token == "" || !new_auth_success {
			api_result := "{\"success\":false, \"message\":\"Unable to get a valid Auth Token from API.\" }"
			return api_result
		} else {
			gomojo_auth_token = new_auth_token
		}
	}

	api_result := ""
	client := &http.Client{}
	api_method := "GET" // overridden later
	var param_data []byte
	var param_reader *bytes.Reader
	param_data = ([]byte)("")

	// Decide on the HTTP Method and Parameters
	if apicall == "auth" {
		api_method = "POST"
		param_data = ([]byte)(apidata)
	} else if apicall == "deauth" {
		api_method = "DELETE"
		apicall = "auth/" + apidata
	} else if apicall == "listoffers" {
		api_method = "GET"
		apicall = "offer"
	} else if apicall == "offerdetails" {
		api_method = "GET"
		apicall = "offer/" + apidata
	} else {
		api_method = "GET"
	}

	// Make the API URL to call
	api_url := "https://www.instamojo.com/api/" + gomojo_api_ver + "/" + apicall + "/"

	param_reader = bytes.NewReader(param_data)

	req, err := http.NewRequest(api_method, api_url, param_reader)
	if err == nil {
		req.Header.Add("X-App-Id", gomojo_app_id)

		if apicall != "auth" {
			req.Header.Add("X-Auth-Token", gomojo_auth_token)
		}

		resp, resperr := client.Do(req)
		if resperr != nil {
			api_result = "{\"success\":false, \"message\":\"Error connecting to or retrieving response from API URL. Please check connectivity. API URL: " + api_url + "\" }"
		} else {
			defer resp.Body.Close()
			bodybytes, _ := ioutil.ReadAll(resp.Body)
			api_result = string(bodybytes)
		}
	}

	return api_result
}
