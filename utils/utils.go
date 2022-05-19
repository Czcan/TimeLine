package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	constSource = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GenerateToken(email string) string {
	h := sha256.New()
	h.Write([]byte(email))
	token := hex.EncodeToString(h.Sum(nil))
	return token
}

func ValidEmail(email string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 {
		return false
	}
	if !regex.MatchString(email) {
		return false
	}
	return true
}

func GenerateNumber(length int, source string) string {
	if source == "" {
		source = constSource
	}
	bytes := []byte(source)
	randNum := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		randNum = append(randNum, bytes[r.Intn(len(bytes))])
	}
	return string(randNum)
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func InferRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
		if Exists(d + "/config.yml") {
			return d
		}
		return infer(filepath.Dir(d))
	}

	return infer(cwd)
}
