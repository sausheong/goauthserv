package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

func main() {
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("goauthserv-secret"))
	m.Use(sessions.Sessions("goauthserv-session", store))

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		IndentJSON: true,
	}))

	// Main page
	m.Get("/", RequireLogin, GetIndex)

	// Login page
	m.Get("/login", GetLogin)

	// Log out the user and redirect the user to the login page
	m.Get("/logout", GetLogout)

	// Authenticate the user
	m.Post("/auth", PostAuth)

	// List of all users
	m.Get("/users", RequireLogin, GetUsers)

	// Add new user page
	m.Get("/users/new", RequireLogin, GetUsersNew)

	// Edit a specific user
	m.Get("/users/:uuid/edit", RequireLogin, GetUsersEdit)

	// Remove an existing user
	m.Get("/users/user/:uuid/remove", RequireLogin, GetUsersRemove)
	//
	// Reset the password for an existing user
	m.Get("/users/user/:uuid/reset", RequireLogin, GetUsersReset)

	// Create a new user or modify an existing user
	m.Post("/users", RequireLogin, PostUsers)

	// APIs from here on

	// Authenticate a user given the user name and password
	m.Post("/authenticate", PostAuthenticate)

	// Validate a session UUID
	m.Post("/validate", PostValidate)

	m.Run()
}
