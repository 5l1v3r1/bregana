package main

import (
  "os"
  "path/filepath"
  "math/rand"
  "time"
  "github.com/pkg/errors"
)


func GetRootPath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "os error")
	}
	dd := filepath.Join(hd, "Bregana")
  os.MkdirAll(dd, 0777)

	return dd, nil
}



func UntestedRandomString(length int) string {
  var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
  const charset = "abcdefghijklmnopqrstuvwxyz1234567890"

  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}
