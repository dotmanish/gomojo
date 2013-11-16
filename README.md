gomojo : Instamojo API Wrapper written in Go
============================================

This is an API wrapper package *and* a command-line tool to interface with 
[Instamojo API](https://www.instamojo.com/developers/), written in Go.

This is currently a work-in-progress, but then, what isn't in this world?

I initially wrote bits of this for Instamojo integration with an app backend.
Now it's a Golang package that can be used to interact with the API.
The latest is at https://github.com/dotmanish/gomojo

    gomojo: is the Go package you can refer to from your Go sources.
    gomojo-tool: is the command-line tool for using the API interactively.

_Disclaimer_: I am not affiliated with Instamojo other than being one of
their  merchant account customers. This (gomojo) is not officially
endorsed by Instamojo team. That said, you should check out
[Instamojo](https://www.instamojo.com/) if you haven't yet.

Install
=======

You would want to do

    go get github.com/dotmanish/gomojo/gomojo-tool

to download and install the *gomojo* package and the *gomojo-tool* binary.
This requires you to have configured GOPATH variable correctly in your
environment.

If you only want to grab the *gomojo* package, just use

    go get github.com/dotmanish/gomojo

which will only retrieve the package, but not the command-line API tool.


Command-Line Tool Usage
=======================

Currently Available actions:

    auth, deauth, listoffers, offerdetails, archiveoffer

Example usage of the command-line API tool:

    gomojo-tool -action listoffers -app <your App-ID> -token <auth token>

    gomojo-tool -action offerdetails -offerslug <offer slug> -app <your App-ID> -token <auth token>

    gomojo-tool -action archiveoffer -offerslug <offer slug> -app <your App-ID> -token <auth token>


If you don't have a pre-generated Auth Token, you can either generate one first like this

    gomojo-tool -action auth -app <your App-ID> -user <your username> -passwd <your password>

Or you can simply pass the username and password as command line parameters
(note that this will generate an intermediate Auth Token, which will be destructed later)
I will later add functionality to read username and password from a file.
 
    gomojo-tool -action listoffers -app <your App-ID> -user <your username> -passwd <your password>

Sample output for 'listoffers':

    Total 2 Offers
    ----------------------------
    Offer # 1
    Status:  Archived
    Title:  Test Product 1
    Slug:  test-product-1
    ShortURL:
    ----------------------------
    Offer # 2
    Status:  Archived
    Title:  Test Product 2
    Slug:  test-product-2
    ShortURL:
    ----------------------------


API Wrapper Package Usage
=========================

The source for *gomojo-tool* is a comprehensive sample for API wrapper usage.
Check *gomojo-tool.go* for reference.

Typical usage would entail:

1. import "github.com/dotmanish/gomojo"

2. call InitGomojoWithAuthToken() or InitGomojoWithUserPass()
   
3. call the Main APIs or Helper Functions
    

**Currently available APIs:**

**Initialization:** 

    InitGomojoWithAuthToken
    InitGomojoWithUserPass

**Main APIs:**

    ListOffers
    GetOfferDetails
    ArchiveOffer
    UploadFile
    CreateOffer
    UpdateOffer
    GetNewAuthToken
    DeleteAuthToken

**Helper Functions:**

    GetCurrentAuthToken


License
=======

Use of this source code is governed by a BSD (3-Clause) License.

Copyright 2013 Manish Malik (manishmalik.name)

All rights reserved.
    
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice,
      this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice,
      this list of conditions and the following disclaimer in the documentation
      and/or other materials provided with the distribution.
    * Neither the name of this program/product nor the names of its contributors may
      be used to endorse or promote products derived from this software without
      specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
