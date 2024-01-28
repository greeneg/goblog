package pbkdf2auth

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"os"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

func readTokensFile(path string) ([]string, error) {
	var records []string

	// now open file and read in tokens
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		records = append(records, scanner.Text())
	}

	return records, scanner.Err()
}

func ValidateViaFile(path string, secret []byte) bool {

	records, err := readTokensFile(path)
	if err != nil {
		return false
	}

	for _, r := range records {
		s := strings.Split(r, ":")

		// make this more clear
		salt, _ := base64.RawStdEncoding.DecodeString(s[2])
		b_salt := []byte(salt)
		var hash []byte = []byte(s[3])

		k := pbkdf2.Key(secret, b_salt, 10000, 20, sha3.New512)

		// now determine if this equals the hash string
		ret := bytes.Compare(k, hash)
		if ret != 0 {
			return false
		}
		break
	}
	return true
}
