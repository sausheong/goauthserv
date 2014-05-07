package utils

import (
  "os"
  "fmt"
  "strings"
  "crypto/sha256"
  "crypto/rand"
  "io/ioutil"
  "code.google.com/p/go.crypto/pbkdf2"
  "github.com/sendgrid/sendgrid-go"
)

// Hashes a plaintext, given a salt and returns it as a string
func Hash(plaintext []byte, salt []byte) string {
  return fmt.Sprintf("%X", pbkdf2.Key(plaintext, salt, 4096, sha256.Size, sha256.New))
}

// Send a reset password
func SendResetPassword(recipient string, password string) {  
  html, _ := ioutil.ReadFile("password_reset.html")
  msg := strings.Replace(string(html), "-password-", password, 1)
  sg := sendgrid.NewSendGridClient(os.Getenv("SGUSER"), os.Getenv("SGPASS"))
  message := sendgrid.NewMail()
  message.AddTo(recipient)
  message.AddToName("goauthserv")
  message.SetSubject("[goauthserv] Password reset")
  message.SetHTML(msg)
  message.SetFrom("noreply@goauthserv")
  if r := sg.Send(message); r != nil {
    fmt.Println(r)
  }  
}

// Randomly generates a password
func RandPassword(str_size int) string {
  alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
  var bytes = make([]byte, str_size)
  rand.Read(bytes)
  for i, b := range bytes {
    bytes[i] = alphanum[b%byte(len(alphanum))]
  }
  return string(bytes)
}
