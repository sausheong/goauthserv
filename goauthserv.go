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
    Directory:  "templates", 
    Layout:     "layout", 
    Extensions: []string{".tmpl", ".html"}, 
    IndentJSON: true, 
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
    session_id, err := gdb.Auth(email, password)    
    fmt.Println("session_id (web):", session_id)
    fmt.Println("err (web):", err)
    if err != nil {
      http.Error(res, "Not Authorized", http.StatusUnauthorized)
    } else {
      
      session.Set("user_session", session_id)
      r.Redirect("/")      
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

