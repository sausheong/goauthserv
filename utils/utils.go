package utils

import (
  "fmt"
  "crypto/sha256"
  "code.google.com/p/go.crypto/pbkdf2"
)

// hashes a plaintext, given a salt and returns it as a string
func Hash(plaintext []byte, salt []byte) string {
  defer clear(plaintext)
  return fmt.Sprintf("%X", pbkdf2.Key(plaintext, salt, 4096, sha256.Size, sha256.New))
}

func clear(b []byte) {
  for i := 0; i < len(b); i++ {
    b[i] = 0;
  }
}
