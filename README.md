gomojo : Instamojo API tool written in Go
=========================================

This is an API wrapper / command-line tool to interface with 
[Instamojo API](https://www.instamojo.com/developers/), written in Go.

This is currently a work-in-progress, but then, what isn't in this world?

I initially wrote bits of this for Instamojo integration with an app backend.
At this stage, this acts like a command-line API tool. The next step in the
process would be to integrate the other bits, and then abstract the
functionality into a Golang package and retain the command-line tool as a 
program which uses the package library.

_Disclaimer_: I am not affiliated with Instamojo other than being one of
their  merchant account customers. This (gomojo) is not officially
endorsed by Instamojo team. That said, you should check out
[Instamojo](https://www.instamojo.com/) if you haven't yet.

Usage
=====

Currently Available actions:

    auth, deauth, listoffers

Example usage of the command-line API tool:

    gomojo-tool -action listoffers -app <your App-ID> -token <auth token>

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


Build
=====

You would want to do

    go build gomojo-tool.go

to create the 'gomojo-tool' or 'gomojo-tool.exe' binary (depending upon your target platform).
If you don't compile and build a binary beforehand, you can replace

    gomojo-tool

with

    go run gomojo-tool.go

in the above usage examples.

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
