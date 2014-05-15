package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/nu7hatch/gouuid"
	"github.com/sausheong/goauthserv/utils"
	"time"
)

type User struct {
	Id        int64
	Uuid      string `sql:"size:255;not null;unique"`
	Email     string `sql:"size:255"`
	Password  string `sql:"size:255"`
	Name      string `sql:"size:255"`
	Salt      string `sql:"size:255"`
  ActivationToken string `sql:"size:255"`
  Activated bool
	CreatedAt time.Time
}

type Session struct {
	Id        int64
	Uuid      string `sql:"size:255;not null;unique"`
	User_id   string
	CreatedAt time.Time
}

type Resource struct {
	Id        int64
	Uuid      string `sql:"size:255;not null;unique"`
	Name      string `sql:"size:255;not null;unique"` // a friendly name for the resource
	CreatedAt time.Time
}

// A permission allows a user to access the resource
type Permission struct {
	Id          int64
	Uuid        string `sql:"size:255;not null;unique"`
	User_id     string
	Resource_id string
	CreatedAt   time.Time
}

var DB gorm.DB

// initialize gorm
func init() {
	var err error
	DB, err = gorm.Open("postgres", "user=goauthserv password=goauthserv dbname=goauthserv sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
	}
}

// format the CreatedAt date to display nicely on the screen
func (u *User) CreatedAtDate() string {
	return u.CreatedAt.Format("01-02-2006")
}

// Before creating a user, add in the uuid
func (u *User) BeforeCreate() (err error) {
	u5, err := uuid.NewV5(uuid.NamespaceURL, []byte(u.Email))
	if err != nil {
		fmt.Println("UUID error:", err)
		return
	}
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Salt error:", err)
		return
	}
	hashed := utils.Hash([]byte(u.Password), []byte(u4.String()))
	u.Password = hashed
	u.Salt = u4.String()
	u.Uuid = u5.String()
	return
}

// Authenticate a user given the user name and the plaintext password
func Auth(email string, password string) (session_id string, err error) {
	// get user from database
	var user = User{}
	err = DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return
	}
	// hash the password
	hashed := utils.Hash([]byte(password), []byte(user.Salt))

	if user.Password == hashed {
		sess := Session{User_id: user.Uuid}
		err = DB.Save(&sess).Error
		if err != nil {
			return
		}
		session_id = sess.Uuid
	} else {
		err = errors.New("Wrong password")
	}
	return
}

// Before creating a session, add in the uuid
func (u *Session) BeforeCreate() (err error) {
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	u.Uuid = u4.String()
	return
}

// Before creating a resource, add in the uuid
func (u *Resource) BeforeCreate() (err error) {
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	u.Uuid = u4.String()
	return
}

// Before creating a permission, add in the uuid
func (u *Permission) BeforeCreate() (err error) {
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	u.Uuid = u4.String()
	return
}

func (u *User) Activate() (err error) {
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
  
  u.ActivationToken = u4.String()
	err = DB.Save(&u).Error
	if err != nil {
		return
	}  
  return
}

func (u *User) Deactivate() (err error) {
  u.ActivationToken = nil
	err = DB.Save(&u).Error
	if err != nil {
		return
	}  
  return  
}
