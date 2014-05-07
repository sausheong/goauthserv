package main

import (
	gdb "goauthserv/db"
)

func main() {
  gdb.DB.AutoMigrate(gdb.User{})
  gdb.DB.AutoMigrate(gdb.Session{})
  user := gdb.User{Name: "Sau Sheong", Email: "sausheong@gmail.com", Password: "123"}
  gdb.DB.Save(&user)
}
