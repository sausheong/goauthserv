[![Build Status](https://travis-ci.org/sausheong/goauthserv.png)](https://travis-ci.org/sausheong/goauthserv)

# goauthserv

goauthserv is a REST-based authentication service written in Go with Martini as the web framework and Gorm as the persistence layer.

## Introduction

goauthserv is a simple REST-based authentication service. Create users in the goauthserv database using its web interface, then allow other applications to access its user database using a simple REST API.

goauthserv answers the following questions:

1. Is this a valid user? (authentication)
2. Is the user currently logged in? (validation)
3. Can the user access this resource? (authorization)

REST API features:

* Authentication of a user using email as account ID (returns a session)
* Validation of a session
* Add a resource (resource is any alphanumeric identifier) 
* Remove a resource 
* Authorization of a user to a resource (NOT IMPLEMENTED YET)
* Check authorization of a user for a resource  (NOT IMPLEMENTED YET)



## Install and setup

Make sure you have at least Go 1.2.1 before you start. Use gvm (https://github.com/moovweb/gvm) -- it's easier.

Create a Postgres database named `goauthserv` with user `goauthserv` and password `goauthserv` with this script:

    ./setupdb
    
Set up the tables with this:
    
    go run setup/setup.go
    
Then install the dependencies like this:

    go get .
    
Now build the server with this:

    go build

Finally run the server with this:

    SGUSER=<SendGrid username> SGPASS=<SendGrid password> PORT=<port number> ./goauthserv
    
Without PORT the server will run at port 3000. A SendGrid account allow you to reset user passwords by sending the user a reset email.

## Clients

A list of example clients (in different languages) using goauthserv can be found in the `clients` folder.

    

