package main

import (
	"crypto/rand"
	"encoding/base64"
)

func GetRandomID() string {
	data := make([]byte, 32)
	rand.Read(data)
	return base64.StdEncoding.EncodeToString(data)
}
