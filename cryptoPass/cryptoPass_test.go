package cryptoPass

import (
	"os"
	"strings"
	"testing"
)

func TestSetSecret(t *testing.T) {
	SetSecret()
	if os.Getenv("HASH_KEY") == "" {
		t.Errorf("HASH_KEY enviorment variable not set")
	}
}

func TestEncryptionCycle(t *testing.T) {
	ciphertext := Encrypt([]byte("password"), os.Getenv("HASH_KEY"))
	plaintext := decrypt(ciphertext, os.Getenv("HASH_KEY"))
	if strings.Compare(string(plaintext), "password") != 0 {
		t.Errorf("decrypted text was : %s it should have been %s", string(plaintext), "password")
	}
}

func TestCompareEncryptedPassAndPass(t *testing.T) {
	SetSecret()
	ciphertext := Encrypt([]byte("password"), os.Getenv("HASH_KEY"))
	if !CompareHashandPassword("password", ciphertext) {
		t.Errorf("The value should have been true but it was false")
	}
	if CompareHashandPassword("notPassword", ciphertext) {
		t.Errorf("The value should have been false but it was true")
	}

}
