package main

import (
 // "fmt"
 "github.com/go-martini/martini" 
 "github.com/martini-contrib/sessions"
 "github.com/martini-contrib/render"
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
  m.Get("/users/edit/:uuid", RequireLogin, GetUsersEdit)
  
  // Remove an existing user
  m.Get("/users/remove/:uuid", RequireLogin, GetUsersRemove)
  // 
  // Reset the password for an existing user
  m.Get("/users/reset/:uuid", RequireLogin, GetUsersReset)
    
  // Create a new user or modify an existing user
  m.Post("/users", RequireLogin, PostUsers)
  // 
  // 
  // //
  // // APIs from here on
  // // 
  // 
  // m.Post("/authenticate", func() string {
  //   return
  // })
  // 
  // m.Post("/validate", func() string {
  //   return
  // })

  m.Run()
}





