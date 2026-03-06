package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32
)

func HashPassword(password string) (string, error) {

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d, t=%d, p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return false, errors.New("unsupported algorithm")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}

	if version != argon2.Version {
		return false, fmt.Errorf("incompatible version: expected %d, got %d", argon2.Version, version)
	}

	var m, t, p uint32
	_, err = fmt.Sscanf(parts[3], "m=%d, t=%d, p=%d", &m, &t, &p)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	expectedHash := argon2.IDKey([]byte(password), salt, t, m, uint8(p), uint32(len(hash)))

	if subtle.ConstantTimeCompare(hash, expectedHash) == 1 {
		return true, nil
	}

	return false, nil
}

func ValidatePassword(password string) error {
	if len(password) < 12 {
		return errors.New("The password must contain at least 12 characters.")
	}

	var (
		hasLowercaseLetter = false
		hasUppercaseLetter = false
		hasNumber          = false
		hasSpecial         = false
	)

	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?/"

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLowercaseLetter = true
		case unicode.IsUpper(char):
			hasUppercaseLetter = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !hasLowercaseLetter {
		return errors.New("The password must contain at least one lowercase letter.")
	}
	if !hasUppercaseLetter {
		return errors.New("The password must contain at least one uppercase letter.")
	}
	if !hasNumber {
		return errors.New("The password must contain at least one number.")
	}
	if !hasSpecial {
		return errors.New("The password must contain at least one special character.")
	}

	return nil
}
