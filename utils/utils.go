package utils

import (
  "fmt"
  "strings"
  "crypto/sha256"
  "crypto/rand"
  "io/ioutil"
  "code.google.com/p/go.crypto/pbkdf2"
  "github.com/sendgrid/sendgrid-go"
)

// hashes a plaintext, given a salt and returns it as a string
func Hash(plaintext []byte, salt []byte) string {
  defer clear(plaintext)
  return fmt.Sprintf("%X", pbkdf2.Key(plaintext, salt, 4096, sha256.Size, sha256.New))
}

func SendResetPassword(recipient string, password string) {  
  html, _ := ioutil.ReadFile("password_reset.html")
  msg := strings.Replace(string(html), "-password-", password, 1)
  sg := sendgrid.NewSendGridClient("sausheong@gmail.com", "chang123")
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

func RandPassword(str_size int) string {
  alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
  var bytes = make([]byte, str_size)
  rand.Read(bytes)
  for i, b := range bytes {
    bytes[i] = alphanum[b%byte(len(alphanum))]
  }
  return string(bytes)
}

func clear(b []byte) {
  for i := 0; i < len(b); i++ {
    b[i] = 0;
  }
}