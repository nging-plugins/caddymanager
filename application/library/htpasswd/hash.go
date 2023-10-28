package htpasswd

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/GehirnInc/crypt/apr1_crypt"
	"golang.org/x/crypto/bcrypt"
)

func hashApr1(password string) (string, error) {
	return apr1_crypt.New().Generate([]byte(password), nil)
}

func hashBcrypt(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ``, err
	}
	return string(passwordBytes), nil
}

func hashSha(password string) string {
	s := sha1.New()
	s.Write([]byte(password))
	passwordSum := []byte(s.Sum(nil))
	return base64.StdEncoding.EncodeToString(passwordSum)
}

type HashAlgorithm string

const (
	// HashAPR1 Apache MD5 crypt - legacy
	HashAPR1 HashAlgorithm = "apr1"
	// HashBCrypt bcrypt - recommended
	HashBCrypt HashAlgorithm = "bcrypt"
	// HashSHA sha5 insecure - do not use
	HashSHA HashAlgorithm = "sha"
)
