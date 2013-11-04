// Copyright 2013 Manish Malik (manishmalik.name)
// All rights reserved.
// Use of this source code is governed by a BSD (3-Clause) License
// that can be found in the LICENSE file.

// I initially wrote bits of this for Instamojo integration with an app backend.
// At this stage, this acts like a command-line API tool. The next step in the
// process would be to integrate the other bits, and then abstract the
// functionality into a Golang package and retain the command-line tool as a
// program which uses the package library.

// Currently Available actions:
//
// auth, deauth, listoffers
//
// Example usage of the command-line API tool:
//
// gomojo-tool -action listoffers -app <your App-ID> -token <auth token>
//
// If you don't have a pre-generated Auth Token, you can either generate one first like this
//
// gomojo-tool -action auth -app <your App-ID> -user <your username> -passwd <your password>
//
// Or you can simply pass the username and password as command line parameters
// (note that this will generate an intermediate Auth Token, which will be destructed later)
// I will later add functionality to read username and password from a file.
//
// gomojo-tool -action listoffers -app <your App-ID> -user <your username> -passwd <your password>
//

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ListOffersResponse: represents response of 'offer' API
type ListOffersResponse struct {
	Offers  []Offer `json:"offers"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
}

// Offer: represents one Offer object (in list of offers)
type Offer struct {
	ShortURL string `json:"shorturl"`
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	Status   string `json:"status"`
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

var cmd_action, cmd_app_id, cmd_auth_token, cmd_username, cmd_passwd, cmd_api_ver string
var authenticated_in_current bool
var gomojo_version string

// Note to people who may read this for learning:
// There are multiple sophisticated command-line option parsers available for Go
// at http://code.google.com/p/go-wiki/wiki/Projects#Command-line_Option_Parsers
// Do take a look at these and decide whether you want to simply use 'flag' package
// or want to instead use the existing available parsing packages.
func init() {

	flag.StringVar(&cmd_action, "action", "", "API Action")
	flag.StringVar(&cmd_app_id, "app", "", "App ID")
	flag.StringVar(&cmd_auth_token, "token", "", "Auth Token")
	flag.StringVar(&cmd_username, "user", "", "Username (for Auth)")
	flag.StringVar(&cmd_passwd, "passwd", "", "Password (for Auth)")
	flag.StringVar(&cmd_api_ver, "version", "1", "API Version (default 1)")

}

// initParams: For command-line usage support
func initParams() {

	var paramsOkay bool
	paramsOkay = true

	gomojo_version = "1.0.2"

	flag.Parse()

	if cmd_action != "auth" && cmd_action != "deauth" && cmd_action != "listoffers" {
		fmt.Print("You must specify the action on command line: 'auth', 'deauth', 'listoffers'\n\n")
		paramsOkay = false
	} else if cmd_app_id == "" {
		fmt.Print("You must specify the App-ID from command line via the '-app' parameter.\n\n")
		paramsOkay = false
	} else if cmd_auth_token == "" && (cmd_username == "" || cmd_passwd == "") {
		fmt.Print("If 'token' is not supplied, then both 'user' and 'password' parameters must be supplied from command line.\n\n")
		paramsOkay = false
	}

	if !paramsOkay {
		fmt.Printf("* gomojo v %s from https://github.com/dotmanish/gomojo\n\n", gomojo_version)
		fmt.Print("Usage: gomojo-tool -action <Action> -app <App IP> [-token <Auth Token>] [-user <Username>] [-passwd <Password>]\n\n")
		fmt.Print("Currently Available actions: auth, deauth, listoffers\n")
		fmt.Print("Example: gomojo-tool -action listoffers -app <your App-ID> -token <auth token>\n")
		fmt.Print("Example: gomojo-tool -action listoffers -app <your App-ID> -user <your username> -passwd <your password>\n")
		os.Exit(1)
	}

}

// processCommandLineAPI: Main handler for command-line usage
// The responsibility of this function is to make sure that
// we display on-screen what's happening and show the API responses
func processCommandLineAPI(apicall string) {

	if apicall == "listoffers" {

		offers, list_success, list_message := listOffers()

		fmt.Println("List Offers API Success:", list_success)
		fmt.Println("List Offers API Message:", list_message)
		fmt.Printf("Total %d Offers\n", len(offers))
		fmt.Println("----------------------------")

		for iter := range offers {
			offer := offers[iter]
			fmt.Println("Offer #", iter+1)
			fmt.Println("Status:", offer.Status)
			fmt.Println("Title:", offer.Title)
			fmt.Println("Slug:", offer.Slug)
			fmt.Println("ShortURL:", offer.ShortURL)
			fmt.Println("----------------------------")
		}

	} else if apicall == "auth" {

		auth_token, auth_success, auth_message := getNewAuthToken(cmd_username, cmd_passwd)

		fmt.Println("New Auth Token:", auth_token)
		fmt.Println("Auth API Success:", auth_success)
		fmt.Println("Auth API Message:", auth_message)

	} else if apicall == "deauth" {

		deauth_success, deauth_message := deleteAuthToken(cmd_auth_token)

		fmt.Println("Delete-Auth API Success:", deauth_success)
		fmt.Println("Delete-Auth API Message:", deauth_message)
	}

}

// listOffers: deletes an existing Auth Token
// Inputs: None
// Returns: (Offer object array, API success bool, Message string)
func listOffers() ([]Offer, bool, string) {

	api_result := callAPI("listoffers", "")

	jsonobj := new(ListOffersResponse)
	jsonerr := json.Unmarshal([]byte(api_result), jsonobj)

	if jsonerr != nil {
		jsonobj.Message = "Invalid JSON: " + jsonerr.Error()
	}

	return jsonobj.Offers, jsonobj.Success, jsonobj.Message
}

// getNewAuthToken: gets a new Auth Token
// Inputs: (Username string, Password string)
// Returns: (Auth Token string, API success bool, Message string)
func getNewAuthToken(username, password string) (string, bool, string) {

	api_result := callAPI("auth", "username="+cmd_username+"&password="+cmd_passwd)

	jsonobj := new(AuthResponse)
	jsonerr := json.Unmarshal([]byte(api_result), jsonobj)

	if jsonerr != nil {
		jsonobj.Message = "Invalid JSON: " + jsonerr.Error()
	}

	return jsonobj.Token, jsonobj.Success, jsonobj.Message
}

// deleteAuthToken: deletes an existing Auth Token
// Inputs: (Auth Token string)
// Returns: (API success bool, Message string)
func deleteAuthToken(auth_token string) (bool, string) {

	api_result := callAPI("deauth", auth_token)

	jsonobj := new(DeAuthResponse)
	jsonerr := json.Unmarshal([]byte(api_result), jsonobj)

	if jsonerr != nil {
		fmt.Println("Invalid JSON: ", jsonerr.Error())
		os.Exit(3)
	}

	return jsonobj.Success, jsonobj.Message
}

// callAPI: Main function handling the REST API
// Valid apicall values: "auth", "deauth", "listoffers"
func callAPI(apicall, apidata string) string {

	// Check if we have auth token available.
	// If not, let's first authenticate and retrieve it.
	if apicall != "auth" && cmd_auth_token == "" {

		cmd_auth_token, _, _ = getNewAuthToken(cmd_username, cmd_passwd)

		if cmd_auth_token != "" {
			// This flag will later tell us if we should delete this auth token
			// since we specifically created this only for current session.
			authenticated_in_current = true
		} else {
			fmt.Print("Failed to get an auth token.\n")
			os.Exit(3)
		}
		// and then continue through to the current API call
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
	} else {
		api_method = "GET"
	}

	// Make the API URL to call
	api_url := "https://www.instamojo.com/api/" + cmd_api_ver + "/" + apicall + "/"

	param_reader = bytes.NewReader(param_data)

	req, err := http.NewRequest(api_method, api_url, param_reader)
	if err == nil {
		req.Header.Add("X-App-Id", cmd_app_id)

		if apicall != "auth" {
			req.Header.Add("X-Auth-Token", cmd_auth_token)
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

func main() {

	initParams()

	processCommandLineAPI(cmd_action)

	if authenticated_in_current {
		processCommandLineAPI("deauth")
	}
}
