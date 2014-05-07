package main

import (
	gdb "github.com/sausheong/goauthserv/db"
)

func main() {
  gdb.DB.Exec("DROP TABLE users;DROP TABLE sessions;")
  gdb.DB.AutoMigrate(gdb.User{})
  gdb.DB.AutoMigrate(gdb.Session{})
  user := gdb.User{Name: "Admin", Email: "admin@goauthserv", Password: "admin"}
  gdb.DB.Save(&user)
}
