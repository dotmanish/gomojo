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
	"fmt"
	"flag"
	"net/http"
	"io/ioutil"
	"os"
	"bytes"
	"encoding/json"
)

type ListOffersResp struct {
	Offers []Offer `json:"offers"`
	Success bool `json:"status"`
}

type Offer struct {
    ShortURL string `json:"shorturl"`
    Title string `json:"title"`
    Slug string `json:"slug"`
    Status string `json:"status"`
}

var cmd_action, cmd_app_id, cmd_auth_token, cmd_username, cmd_passwd, cmd_api_ver string
var authenticated_in_current bool

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

func initParams() {

	var paramsOkay bool
    paramsOkay = true

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
                fmt.Print("Usage: gomojo-tool -action <Action> -app <App IP> [-token <Auth Token>] [-user <Username>] [-passwd <Password>]\n\n")
                fmt.Print("Currently Available actions: auth, deauth, listoffers\n")
                fmt.Print("Example: gomojo-tool -action listoffers -app <your App-ID> -token <auth token>\n")
                fmt.Print("Example: gomojo-tool -action listoffers -app <your App-ID> -user <your username> -passwd <your password>\n")
                os.Exit(1)
        }

}

func filterOutputAPI(apicall string) {

	api_result := callAPI(apicall)
	
	if apicall == "listoffers" {

		jsonobj := new(ListOffersResp)
    	jsonerr := json.Unmarshal([]byte(api_result), jsonobj)
		
		if jsonerr != nil {
			fmt.Println("Invalid JSON: ",jsonerr.Error())
			os.Exit(3)
		} else {
			fmt.Printf("Total %d Offers\n",len(jsonobj.Offers))
			fmt.Println("----------------------------")

			for iter := range jsonobj.Offers {
				offer := jsonobj.Offers[iter]
				fmt.Println("Offer #",iter+1)
				fmt.Println("Status: ",offer.Status)
				fmt.Println("Title: ",offer.Title)
				fmt.Println("Slug: ",offer.Slug)
				fmt.Println("ShortURL: ",offer.ShortURL)
				fmt.Println("----------------------------")
			}	
		}
		
	} else {
		fmt.Println(api_result)
	}

}

// Valid apicall values:
// "auth", "deauth", "listoffers"
func callAPI(apicall string) string {

	// Check if we have auth token available.
	// If not, let's first authenticate and retrieve it.
	if apicall != "auth" && cmd_auth_token == "" {
		cmd_auth_token = callAPI("auth")
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
	client := &http.Client{ }
	api_method := "GET" // overridden later
	var param_data []byte
	var param_reader *bytes.Reader
	param_data = ([]byte)("")

	// Decide on the HTTP Method and Parameters
	if apicall == "auth" {
		api_method = "POST"
		param_data = ([]byte)("username=" + cmd_username + "&password=" + cmd_passwd + "&debug=true")
	} else if apicall == "deauth" {
		api_method = "DELETE"
		apicall = "auth/:" + cmd_auth_token
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
		
		if (apicall != "auth") {
			req.Header.Add("X-Auth-Token", cmd_auth_token)
		}
		
		resp, resperr := client.Do(req)
		if resperr != nil {
			fmt.Print("Error connecting to or retrieving response from API URL. Please check connectivity.\n")
			fmt.Print("API URL: " + api_url + "\n")
			os.Exit(2)
		}

		defer resp.Body.Close()
	
		bodybytes, _ := ioutil.ReadAll(resp.Body)

		api_result = string(bodybytes)
	}

	return api_result;
}

func main() {

	initParams()
	
	filterOutputAPI(cmd_action)

	if authenticated_in_current {
		filterOutputAPI("deauth")
	}
}
