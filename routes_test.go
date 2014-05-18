package main

import (
	gdb "github.com/sausheong/goauthserv/db"

	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
  // "github.com/martini-contrib/sessions"
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

func Test_GetUsersEdit(t *testing.T) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/users/user/:uuid/edit", GetUsersEdit)

	user := create_user("Sau Sheong", "sausheong@me.com", "123")

	defer delete_user(user)
	res := httptest.NewRecorder()
	url := strings.Join([]string{"/users/user/", user.Uuid, "/edit"}, "")
	req, _ := http.NewRequest("GET", url, nil)

	m.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Response code is %v", res.Code)
	}
}

func Test_GetUsersRemove(t *testing.T) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/users/user/:uuid/remove", GetUsersRemove)

	user := create_user("Sau Sheong", "sausheong@me.com", "123")

	defer delete_user(user)
	res := httptest.NewRecorder()
	url := strings.Join([]string{"/users/user/", user.Uuid, "/remove"}, "")
	req, _ := http.NewRequest("GET", url, nil)

	m.ServeHTTP(res, req)
	if res.Code != 302 {
		t.Errorf("Response code is %v", res.Code)
	}
}

func Test_GetUsersActivate(t *testing.T) {
  m := martini.Classic()
  m.Use(render.Renderer())
  user := create_user("Sau Sheong", "sausheong@me.com", "123")
  defer delete_user(user)
  
  activation_url, err := user.ActivationUrl()
  if err != nil {
    t.Errorf("Error is %v", err)
  }
  m.Get("/users/:uuid/activate", GetUsersActivate)
  

  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", activation_url, nil)

  m.ServeHTTP(res, req)
  if res.Code != 200 {
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
	var json_msg interface{}
	if err := json.Unmarshal(to_byte_array(res.Body), &json_msg); err != nil {
		t.Errorf("Cannot get JSON, error is %v", err)
	}
}

func Test_PostValidate(t *testing.T) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Post("/authenticate", PostAuthenticate)
	m.Post("/validate", PostValidate)

	user := create_user("Sau Sheong", "sausheong@me.com", "123")
	defer delete_user(user)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/authenticate", nil)
	req.ParseForm()
	req.PostForm.Add("email", "sausheong@me.com")
	req.PostForm.Add("password", "123")

	m.ServeHTTP(res, req)
	var json_msg interface{}
	if err := json.Unmarshal(to_byte_array(res.Body), &json_msg); err != nil {
		t.Errorf("Cannot get JSON, error is %v", err)
	}
	msg := json_msg.(map[string]interface{})

	req, _ = http.NewRequest("POST", "/validate", nil)
	req.ParseForm()
	req.PostForm.Add("session", msg["session"].(string))

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

func to_byte_array(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func create_user(name string, email string, password string) gdb.User {
	user := gdb.User{Name: name, Email: email, Password: password}
	err := gdb.DB.DB().Ping()
	if err != nil {
		return user
	} else {
		gdb.DB.Save(&user)
		return user
	}
}

func delete_user(user gdb.User) (err error) {
	db_err := gdb.DB.DB().Ping()
	if db_err != nil {
		return
	} else {
		err = gdb.DB.Delete(&user).Error
		return
	}
}
