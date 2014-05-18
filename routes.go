package main

import (
	"net/http"
	"goauthserv/db"
  
  "github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"	
	"github.com/sausheong/goauthserv/utils"

)

// GET /
func GetIndex(r render.Render) {
	r.HTML(200, "index", nil)
}

// GET /login
func GetLogin(r render.Render) {
	r.HTML(200, "login", nil, render.HTMLOptions{Layout: ""}) // no layout
}

// GET /logout
func GetLogout(s sessions.Session, r render.Render) {
	s.Clear()
	r.Redirect("/login")
}

// POST /auth
func PostAuth(s sessions.Session, r render.Render, req *http.Request) {
	email := req.PostFormValue("email")
	password := req.PostFormValue("password")
	session_id, err := db.Auth(email, password)
	if err != nil {
		r.Error(401)
	} else {
		s.Set("user_session", session_id)
		r.Redirect("/")
	}
}

// GET /users
func GetUsers(r render.Render) {
	users := []db.User{}
	if db.DB.Find(&users).RecordNotFound() {
		r.Error(404)
	} else {
		r.HTML(200, "users", users)
	}
}

// GET /users/new
func GetUsersNew(r render.Render) {
	r.HTML(200, "users/new", nil)
}

// GET /users/user/:uuid/edit
func GetUsersEdit(r render.Render, params martini.Params) {
	user := db.User{}
	if db.DB.Where("uuid = ?", params["uuid"]).First(&user).RecordNotFound() {
		r.Error(404)
	} else {
		r.HTML(200, "users/edit", user)
	}
}

// GET /users/user/:uuid/remove
func GetUsersRemove(r render.Render, params martini.Params) {
	user := db.User{}
	if db.DB.Where("uuid = ?", params["uuid"]).First(&user).RecordNotFound() {
		r.Error(404)
	} else {
		if err := db.DB.Delete(&user).Error; err != nil {
			r.Error(500)
		} else {
			r.Redirect("/users")
		}
	}
}

// GET /users/user/:uuid/reset
func GetUsersReset(r render.Render, params martini.Params) {
	user := db.User{}
	if db.DB.Where("uuid = ?", params["uuid"]).First(&user).RecordNotFound() {
		r.Error(404)
	} else {
		password := utils.RandPassword(8)
		user.Password = utils.Hash([]byte(password), []byte(user.Salt))
		db.DB.Save(&user)
		go utils.SendResetPassword(user.Email, password)
		r.Redirect("/users")
	}
}

// GET /users/:uuid/activate
func GetUsersActivate(r render.Render, params martini.Params) {
  user := db.User{}
  if db.DB.Where("activation_token = ?", params["uuid"]).First(&user).RecordNotFound() {
    r.Error(404)
  } else {
    if err := user.Activate(); err != nil {
      r.Error(500)
    }
    r.Status(200)
  }
}

// POST /users
func PostUsers(r render.Render, req *http.Request) {
	name := req.PostFormValue("name")
	email := req.PostFormValue("email")
	password := req.PostFormValue("password")
	uuid := req.PostFormValue("uuid")
	var user = db.User{}
	if uuid != "" {
		if db.DB.Where("uuid = ?", uuid).First(&user).RecordNotFound() {
			r.Error(404)
		}
		user.Name = name
		user.Email = email
	} else {
		user = db.User{Name: name, Email: email, Password: password}
	}
	if err := db.DB.Save(&user).Error; err != nil {
		r.Error(500)
	} else {
		r.Redirect("/users")
	}
}

// POST /authenticate
func PostAuthenticate(r render.Render, req *http.Request) {
	email := req.PostFormValue("email")
	password := req.PostFormValue("password")
	session_id, err := db.Auth(email, password)
	if err != nil {
		r.Error(401)
	} else {
		r.JSON(200, map[string]interface{}{"session": session_id})
	}
}

// POST /validate
func PostValidate(r render.Render, req *http.Request) {
	s := req.PostFormValue("session")
	session := db.Session{}
	if db.DB.Where("uuid = ?", s).First(&session).RecordNotFound() {
		r.Error(404)
	} else {
		r.Status(200)
	}
}

// POST /authorize
// Authorize a user to use a resource
func PostAuthorize(r render.Render, req *http.Request) {

}

// POST /check
// Check if a user is authorized to use a resource
func PostCheck(r render.Render, req *http.Request) {

}

// POST /resource/add
// Add a resource
func PostResourceAdd(r render.Render, req *http.Request) {

}

// POST /resource/remove
// Remove a resource
func PostResourceRemove(r render.Render, req *http.Request) {

}

// Handler to require a user to log in. If the user is currently logged in
// nothing happens. Otherwise clear existing session and redirect the user
// to the login page
func RequireLogin(s sessions.Session, r render.Render) {
	session := s.Get("user_session")
	if session == nil {
		s.Clear()
		r.Redirect("/login")
	}
}
