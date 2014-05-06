package main

import (
 gdb "goauthserv/db"
 // "fmt"
 "net/http"
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
  m.Get("/", require_login, func(s sessions.Session, r render.Render, res http.ResponseWriter) {        
    r.HTML(200, "index", nil)
  })
  
  // Login page
  m.Get("/login", func(r render.Render) {
    r.HTML(200, "login", nil, render.HTMLOptions{Layout: ""}) // no layout
  })
  
  // Log out the user and redirect the user to the login page
  m.Get("/logout", func(s sessions.Session, r render.Render) {
    s.Clear()
    r.Redirect("/login")
  })
  
  // Authenticate the user
  m.Post("/auth", func(r render.Render, req *http.Request, res http.ResponseWriter, session sessions.Session) {
    email := req.PostFormValue("email")
    password := req.PostFormValue("password")

    session_id, err := gdb.Auth(email, password)    
    if err != nil {
      http.Error(res, "Not Authorized", http.StatusUnauthorized)
    } else {      
      session.Set("user_session", session_id)
      r.Redirect("/")      
    }
  })
  
  // List of all users
  m.Get("/users", require_login, func(r render.Render)  {
    users := []gdb.User{}
    gdb.DB.Find(&users)
    r.HTML(200, "users", users)
  })
  
  // Add new user page
  m.Get("/users/new", require_login, func(r render.Render)  {
    r.HTML(200, "users.new", nil)
  })
  
  // Edit a specific user
  m.Get("/users/edit/:uuid", func(r render.Render, params martini.Params) {
    user := gdb.User{}
    gdb.DB.Where("uuid = ?", params["uuid"]).First(&user)
    r.HTML(200, "users.edit", user)
  })
  
  // Remove an existing user
  // m.Get("/users/remove/:uuid", func(params martini.Params) string {
  //   return
  // })
  // 
  // Reset the password for an existing user
  // m.Get("/users/reset/:uuid", func(params martini.Params) string {
  //   return
  // })
  // 
  
  // Create a new user or modify an existing user
  m.Post("/users", func(r render.Render, req *http.Request, params martini.Params) {
    name := req.PostFormValue("name")
    email := req.PostFormValue("email") 
    uuid := req.PostFormValue("uuid")
    var user = gdb.User{}
    if uuid != "" {
      gdb.DB.Where("uuid = ?", uuid).First(&user)
      user.Name = name
      user.Email = email
    } else {
      user = gdb.User{Name: name, Email: email}      
    }
    gdb.DB.Save(&user) 
    r.Redirect("/users")    
  })
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


// Handler to require a user to log in. If the user is currently logged in
// nothing happens. Otherwise clear existing session and redirect the user 
// to the login page
func require_login(sess sessions.Session, r render.Render) {
  s := sess.Get("user_session")
  if  s == nil {
    s.Clear()
    r.Redirect("/login")
  }
}


