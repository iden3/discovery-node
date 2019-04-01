package discovery

import (
	"crypto/rand"
	"encoding/base64"
)

func randBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func randStr(n int) (string, error) {
	b, err := randBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
