package main

import (
  gdb "github.com/sausheong/goauthserv/db"
  
	"net/http"
	"net/http/httptest"
	"testing"
  "io"
  "bytes"
  "io/ioutil"  
  
  "github.com/go-martini/martini" 
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/sessions"
)

func Test_GetIndex(t *testing.T) {
  m := martini.Classic()
  m.Use(render.Renderer())
  m.Get("/", GetIndex)
  
  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/", nil)
  
  m.ServeHTTP(res, req)
  
  if res.Code != 200 {
    t.Errorf("Response code is %v", res.Code)
  }
  html, _ := ioutil.ReadFile("templates/index.tmpl")
  if to_string(res.Body) != string(html) {
    t.Errorf("Expected %v but got his %v", string(html), res.Body)
  }
}

func Test_GetLogin(t *testing.T) {
  m := martini.Classic()
  m.Use(render.Renderer())
  m.Get("/login", GetLogin)
  
  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/login", nil)
  
  m.ServeHTTP(res, req)
  
  if res.Code != 200 {
    t.Errorf("Response code is %v", res.Code)
  }
}

func Test_GetUsers(t *testing.T) {
  m := martini.Classic()
  m.Use(render.Renderer())
  m.Get("/users", GetUsers)
  
  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/users", nil)
  
  m.ServeHTTP(res, req)
  
  if res.Code != 200 {
    t.Errorf("Response code is %v", res.Code)
  }  
}

func Test_GetUsersNew(t *testing.T) {
  m := martini.Classic()
  m.Use(render.Renderer())
  m.Get("/users/new", GetUsersNew)
  
  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/users/new", nil)
  
  m.ServeHTTP(res, req)
  
  if res.Code != 200 {
    t.Errorf("Response code is %v", res.Code)
  }  
}

func Test_PostAuth(t *testing.T) {
  m := martini.Classic()
  m.Use(render.Renderer())
  store := sessions.NewCookieStore([]byte("test-goauthserv"))
  m.Use(sessions.Sessions("test_session", store))
  m.Post("/auth", PostAuth)
  
  user := create_user("Sau Sheong", "sausheong@me.com", "123")
  defer delete_user(user)
  res := httptest.NewRecorder()
  req, _ := http.NewRequest("POST", "/auth", nil)
  req.ParseForm()
  req.PostForm.Add("email", "sausheong@me.com")
  req.PostForm.Add("password", "123")
  
  m.ServeHTTP(res, req)  
  if res.Code != 302 {
    t.Errorf("Response code is %v", res.Code)
  }
  
}

func Test_PostAuthenticate(t *testing.T) {
  m := martini.Classic()
  m.Use(render.Renderer())
  m.Post("/authenticate", PostAuthenticate)
  
  user := create_user("Sau Sheong", "sausheong@me.com", "123")
  defer delete_user(user)
  res := httptest.NewRecorder()
  req, _ := http.NewRequest("POST", "/authenticate", nil)
  req.ParseForm()
  req.PostForm.Add("email", "sausheong@me.com")
  req.PostForm.Add("password", "123")
  
  m.ServeHTTP(res, req)  
  if res.Code != 200 {
    t.Errorf("Response code is %v", res.Code)
  }
  
}

// test helper functions

func to_string(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}

func create_user(name string, email string, password string) gdb.User {
  user := gdb.User{Name: name, Email: email, Password: password}      
  gdb.DB.Save(&user)  
  return user
}

func delete_user(user gdb.User) (err error) {
  err = gdb.DB.Delete(&user).Error
  return
}