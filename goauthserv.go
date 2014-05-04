package main

import (
 gdb "goauthserv/db"
 "fmt"
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
    Directory: "templates", // Specify what path to load the templates from.
    Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
    Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
    Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
    IndentJSON: true, // Output human readable JSON
  }))
  
  m.Get("/", func(r render.Render) {    
    r.HTML(200, "index", nil)
  })
  
  
  m.Get("/login", func(r render.Render) {
    r.HTML(200, "login", nil, render.HTMLOptions{Layout: ""}) // no layout
  })
  
  // m.Get("/logout", func(r render.Render) string {
  //   return
  // })
  // 
  // 
  m.Post("/auth", func(r render.Render, req *http.Request, res http.ResponseWriter, session sessions.Session) {
    err := req.ParseForm()
    if err != nil {
      http.Error(res, "Not Authorized", http.StatusUnauthorized)
    }    
    email := req.PostFormValue("email")
    password := req.PostFormValue("password")
    
    fmt.Println("email:", email)
    fmt.Println("password:", password)
    auth := gdb.Auth(email, password)    

    if auth {
      session.Set("user_session", )
      r.Redirect("/")
    } else {
      http.Error(res, "Not Authorized", http.StatusUnauthorized)
    }
    
  })
  // 
  // 
  // m.Get("/users", func() string {
  //   return
  // })
  // 
  // 
  // m.Get("/users/user/:uuid", func(params martini.Params) string {
  //   return
  // })
  // 
  // m.Get("/users/new", func() string {
  //   return
  // })
  // 
  // m.Get("/users/edit/:uuid", func(params martini.Params) string {
  //   return
  // })
  // 
  // 
  // m.Get("/users/remove/:uuid", func(params martini.Params) string {
  //   return
  // })
  // 
  // 
  // m.Get("/users/reset/:uuid", func(params martini.Params) string {
  //   return
  // })
  // 
  m.Post("/users", func(params martini.Params) {
    user := gdb.User{Name: "sausheong"}
    gdb.DB.Save(&user)
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

