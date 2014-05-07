# goauthserv

goauthserv is a REST-based authentication service written in Go with Martini as the web framework and Gorm as the persistence layer.

## Setup

Make sure you have at least Go 1.2.1.

Create a Postgres database named `goauthserv` with user `goauthserv` and password `goauthserv`. Set up the tables with this:
    
    go run setup/setup.go
    
Then install the dependencies like this:

    go get .
    
Now build the server with this:

    go build

Finally run the server with this:

    ./goauthserv


    

