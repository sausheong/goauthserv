package db

import (
  "fmt"
  "time"
  "github.com/nu7hatch/gouuid"  
  "github.com/jinzhu/gorm"
  _ "github.com/lib/pq"
  gutil "goauthserv/utils"
)

type User struct {
  Id           int64
  Uuid         string  `sql:"size:255;not null;unique"`
  Email        string  `sql:"size:255"`
  Password     string  `sql:"size:255"`
  Name         string  `sql:"size:255"`
  Salt         string  `sql:"size:255"`
  CreatedAt    time.Time
}

type Session struct {
  Id           int64
  Uuid         string  `sql:"size:255;not null;unique"`
  User_id      int64
  CreatedAt    time.Time
}

var DB gorm.DB
func init() {
  var err error
  DB, err = gorm.Open("postgres", "user=goauthserv password=goauthserv dbname=goauthserv sslmode=disable")
  if err != nil {
      panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
  }
}

// Before creating a user, add in the uuid
func (u *User) BeforeCreate() (err error) {
  u5, err := uuid.NewV5(uuid.NamespaceURL, []byte("goauthserv-user"))
  if err != nil {
      fmt.Println("UUID error:", err)
      return
  }
  u4, err := uuid.NewV4()
  if err != nil {
      fmt.Println("Salt error:", err)
      return
  }
  fmt.Println("user email", u.Email)
  fmt.Println("user password", u.Password)
  
  hashed := gutil.Hash([]byte(u.Password), []byte(u4.String()))  
  
  u.Password = hashed  
  u.Salt = u4.String()
  u.Uuid = u5.String()
  return
}

// Authenticate a user given the user name and the plaintext password
// returns a http.HandlerFunc
func Auth(username string, password string) http.HandlerFunc {
  // get user from database
  var user = User{} 
  DB.Where("email = ?", username).First(&user)  
  // hash the password
  hashed := gutil.Hash([]byte(password), []byte(user.Salt))  
  // return a function thatrejects the user as unauthorized if the hashed password doesn't match the one in the database
  return func(res http.ResponseWriter, req *http.Request) {
    if user.Password != hashed {
      http.Error(res, "Not Authorized", http.StatusUnauthorized)
    }
  }
}

// Authenticate a user given the user name and the plaintext password
func Auth(email string, password string) bool {
  // get user from database
  var user = User{} 
  DB.Where("email = ?", email).First(&user)  
  // hash the password
  fmt.Println("user password:", user.Password)
  hashed := gutil.Hash([]byte(password), []byte(user.Salt))  
  
  fmt.Println("hashed password:", hashed)
  if user.Password == hashed {
    return true
  } else {
    return false
  }
}

// Before creating a session, add in the uuid
func (u *Session) BeforeCreate() (err error) {
  u5, err := uuid.NewV5(uuid.NamespaceURL, []byte("goauthserv-session"))
  if err != nil {
      fmt.Println("error:", err)
      return
  }
  u.Uuid = u5.String()
  return
}