// Copyright 2013 Manish Malik (manishmalik.name)
// All rights reserved.
// Use of this source code is governed by a BSD (3-Clause) License
// that can be found in the LICENSE file.

// I initially wrote bits of this for Instamojo integration with an app backend.
// Now it's a Golang package that can be used to interact with the API.
// The latest is at https://github.com/dotmanish/gomojo
//
// This is the Command-Line Tool that uses the API wrapper package (gomojo)
//
// Currently Available actions:
//
// auth, deauth, listoffers, offerdetails
//
// Example usage of the command-line API tool:
//
// gomojo-tool -action listoffers -app <your App-ID> -token <auth token>
//
// gomojo-tool -action offerdetails -offerslug <offer slug> -app <your App-ID> -token <auth token>
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
	"flag"
	"fmt"
	"os"

	"github.com/dotmanish/gomojo"
)

var cmd_action, cmd_app_id, cmd_auth_token, cmd_username, cmd_passwd, cmd_api_ver string
var cmd_offer_slug string
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
	flag.StringVar(&cmd_offer_slug, "offerslug", "", "Offer Slug")
	flag.StringVar(&cmd_api_ver, "version", "1", "API Version (default 1)")

}

// initParams: For command-line usage support
func initParams() {

	var paramsOkay bool
	paramsOkay = true

	gomojo_version := "1.0.2"

	flag.Parse()

	if cmd_action != "auth" && cmd_action != "deauth" && cmd_action != "listoffers" && cmd_action != "offerdetails" {
		fmt.Print("You must specify the action on command line: 'auth', 'deauth', 'listoffers', 'offerdetails'\n\n")
		paramsOkay = false
	} else if cmd_app_id == "" {
		fmt.Print("You must specify the App-ID from command line via the '-app' parameter.\n\n")
		paramsOkay = false
	} else if cmd_auth_token == "" && (cmd_username == "" || cmd_passwd == "") {
		fmt.Print("If 'token' is not supplied, then both 'user' and 'password' parameters must be supplied from command line.\n\n")
		paramsOkay = false
	} else if cmd_action == "offerdetails" && cmd_offer_slug == "" {
		fmt.Print("You must specifiy the Offer Slug via the command line option -offerslug to get offer details.\n\n")
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

		offers, list_success, list_message := gomojo.ListOffers()

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

	} else if apicall == "offerdetails" {

		offer, details_success, details_message := gomojo.GetOfferDetails(cmd_offer_slug)

		fmt.Println("Offer Details API Success:", details_success)
		fmt.Println("Offer Details API Message:", details_message)

		if details_success {
			fmt.Println("----------------------------")
			fmt.Println("Status:", offer.Status)
			fmt.Println("Title:", offer.Title)
			fmt.Println("Slug:", offer.Slug)
			fmt.Println("ShortURL:", offer.ShortURL)
			fmt.Println("Base Price:", offer.BasePrice)
			fmt.Println("Currency:", offer.Currency)
			fmt.Println("Quantity:", offer.Quantity)
			fmt.Println("Start Date:", offer.StartDate)
			fmt.Println("End Date:", offer.EndDate)
			fmt.Println("Timezone:", offer.Timezone)
			fmt.Println("Venue:", offer.Venue)
			fmt.Println("RedirectURL:", offer.RedirectURL)
			fmt.Println("Note:", offer.Note)
			fmt.Println("Description:", offer.Description)
			fmt.Println("----------------------------")
		}

	} else if apicall == "auth" {

		auth_token, auth_success, auth_message := gomojo.GetNewAuthToken(cmd_username, cmd_passwd)

		fmt.Println("New Auth Token:", auth_token)
		fmt.Println("Auth API Success:", auth_success)
		fmt.Println("Auth API Message:", auth_message)

		gomojo.SetCurrentAuthToken(auth_token)

	} else if apicall == "deauth" {

		deauth_success, deauth_message := gomojo.DeleteAuthToken(cmd_auth_token)

		fmt.Println("Delete-Auth API Success:", deauth_success)
		fmt.Println("Delete-Auth API Message:", deauth_message)
	}

}

func main() {

	initParams()

	// Decide how to initialize gomojo
	if cmd_auth_token != "" {
		gomojo.InitGomojoWithAuthToken(cmd_api_ver, cmd_app_id, cmd_auth_token)
	} else {
		gomojo.InitGomojoWithUserPass(cmd_api_ver, cmd_app_id, cmd_username, cmd_passwd)

		if cmd_action != "auth" {
			authenticated_in_current = true
		}
	}

	processCommandLineAPI(cmd_action)

	// Destuct any temporary Auth Tokens we generated specifically for this session
	// (except when 'auth' command-line action was specified).
	if authenticated_in_current {
		fmt.Println("Destructing the Auth Token generated specifically for this session.")
		cmd_auth_token = gomojo.GetCurrentAuthToken()
		processCommandLineAPI("deauth")
	}
}
