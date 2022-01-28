package utils

import (
	"encoding/base64"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func Compare(hash, s string) bool {
	incoming := []byte(s)
	existing := []byte(hash)
	compare := bcrypt.CompareHashAndPassword(existing, incoming)
	return compare == nil
}

func ToBase64(b []byte) string {
	var base64Encoding string

	file := base64.StdEncoding.EncodeToString(b)

	// Determine the content type of the image file
	mimeType := http.DetectContentType(b)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += file

	return base64Encoding
}

func Base64ToByte(b string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b)
}

func FromBase64(b string) (string, string) {
	var typeFile string
	if strings.Contains(b, "jpeg") {
		typeFile = "jpeg"
	} else {
		typeFile = "png"
	}

	return typeFile, b[strings.IndexByte(b, ',')+1:]
}
