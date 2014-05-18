package main

import (
	"goauthserv/db"
)

func main() {
  db.DB.Exec("DROP TABLE users;DROP TABLE sessions;DROP TABLE resources; DROP TABLE permissions;")
  db.DB.AutoMigrate(db.User{})
  db.DB.AutoMigrate(db.Session{})
  db.DB.AutoMigrate(db.Resource{})
  db.DB.AutoMigrate(db.Permission{})
  user := db.User{Name: "Admin", Email: "admin@goauthserv", Password: "admin"}
  db.DB.Save(&user)
}
