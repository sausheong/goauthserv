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
  
  m.Get("/", require_login, func(s sessions.Session, r render.Render, res http.ResponseWriter) {        
    r.HTML(200, "index", nil)
  })
  
  
  m.Get("/login", func(r render.Render) {
    r.HTML(200, "login", nil, render.HTMLOptions{Layout: ""}) // no layout
  })
  
  m.Get("/logout", func(s sessions.Session, r render.Render) {
    s.Clear()
    r.Redirect("/login")
  })
  
  
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
  
  
  m.Get("/users", require_login, func(r render.Render)  {
    users := []gdb.User{}
    gdb.DB.Find(&users)
    r.HTML(200, "users", users)
  })
  
  
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

// Authenticate a user given the user name and the plaintext password
// returns a http.HandlerFunc

func require_login(sess sessions.Session, r render.Render) {
  s := sess.Get("user_session")
  if  s == nil {
    r.Redirect("/login")
  }

}


