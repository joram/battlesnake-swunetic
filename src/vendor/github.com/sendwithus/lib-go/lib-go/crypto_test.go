package swu

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestEncryptAES(t *testing.T) {
	os.Setenv("SWU_CRYPTO_AES_PRIVATE_KEY", "s,ad98!@#a[]asdm")
	original_bytes := []byte("decrypted string!")
	padding_byte := byte(' ')

	encrypted, public_iv, err := EncryptAES(&original_bytes, padding_byte)
	if err != nil {
		t.Error(err)
	}

	decrypted, err := DecryptAES(encrypted, public_iv, padding_byte)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(original_bytes, *decrypted) {
		t.Error("Expected %s, got %s ", original_bytes, *decrypted)
	}
}

func TestEncryptAESEmptyString(t *testing.T) {
	os.Setenv("SWU_CRYPTO_AES_PRIVATE_KEY", "s,ad98!@#a[]asdm")
	original_bytes := []byte("")
	padding_byte := byte(' ')

	encrypted, public_iv, err := EncryptAES(&original_bytes, padding_byte)
	if err != nil {
		t.Error(err)
	}

	decrypted, err := DecryptAES(encrypted, public_iv, padding_byte)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(original_bytes, *decrypted) {
		t.Error("Expected %s, got %s ", original_bytes, decrypted)
	}
}

func TestDecryptAESNoEnv(t *testing.T) {
	os.Unsetenv("SWU_CRYPTO_AES_PRIVATE_KEY")

	encrypted_json_string := []byte("i am a jason")
	byte_public_key := []byte("I am an iv")
	padding_byte := byte(' ')

	_, err := DecryptAES(&encrypted_json_string, &byte_public_key, padding_byte)

	if !strings.Contains(err.Error(), "Missing AES") {
		t.Error("Did not return error for nil private_key")
	}
}
