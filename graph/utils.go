package graph

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/matthewhartstonge/argon2"
	"github.com/spf13/cast"
)

func EncodeID(typeName string, id uint) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", typeName, id)))
}

func DecodeID(id string) (string, uint, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return "", 0, err
	}
	decodedString := string(decodedBytes)
	splitString := strings.Split(decodedString, ":")
	if len(splitString) >= 2 {
		typeName := splitString[0]
		id, err := strconv.ParseUint(splitString[1], 10, 32)
		if err != nil {
			return "", 0, err
		}
		uid, err := cast.ToUintE(id)
		if err != nil {
			return "", 0, err
		}
		return typeName, uid, nil
	}
	return "", 0, errors.New("invalid id")
}

func EncodePassword(password string) ([]byte, error) {
	argon := argon2.DefaultConfig()

	return argon.HashEncoded([]byte(password))
}

func VerifyPassword(encoded []byte, password string) bool {
	ok, err := argon2.VerifyEncoded([]byte(password), encoded)
	if err != nil {
		return false
	}

	return ok
}
