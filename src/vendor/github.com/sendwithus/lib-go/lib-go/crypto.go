package swu

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"os"
)

func padSlice(src *[]byte, block_size int, pad_char byte) *[]byte {
	// Multiple
	mult := int((len(*src) / block_size) + 1)
	// Total length for next block size multiple
	length := block_size * mult
	// Make []byte slice of the correct Block length to copy original []bytes into
	src_padded := make([]byte, length)
	// Get number of elements copied into it. length - elements_copied will be the number
	// of bytes we need to pad
	elements_copied := copy(src_padded, *src)
	// Add padding char to the rest of the slice
	for i := elements_copied; i < length; i++ {
		src_padded[i] = pad_char
	}

	return &src_padded
}

// EncryptAES encrypts a byte array using AES CBC
// The decrypted value is padded with padding_value before encrypting
// It returns decrypted, public_iv *[]bytes and any error encountered
func EncryptAES(decrypted_bytes *[]byte, padding_value byte) (encrypted *[]byte, public_key *[]byte, err error) {
	// Get AES key from environment and return err if it was not present
	private_key := []byte(os.Getenv("SWU_CRYPTO_AES_PRIVATE_KEY"))
	if len(private_key) == 0 {
		return nil, nil, errors.New("Missing AES private key")
	}

	// Make random byte 16 bit iv
	iv := make([]byte, 16)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, nil, err
	}

	// Create new aes block to be passed to the Decrypter
	cipher_block, err := aes.NewCipher(private_key)
	if err != nil {
		return nil, nil, err
	}

	// Pad it with padding_byte so that CryptBlocks works, it only takes full blocks of 16
	decrypted_bytes = padSlice(decrypted_bytes, len(private_key), padding_value)

	// Encrypt!
	mode := cipher.NewCBCEncrypter(cipher_block, iv)
	encrypted_bytes := make([]byte, len(*decrypted_bytes))
	mode.CryptBlocks(encrypted_bytes, *decrypted_bytes)

	return &encrypted_bytes, &iv, err
}

// DecryptAES decrypts a AES CBC encrypted string
// padding_value is trimmed from the decrypted value before returning
// It returns the decrypted *[]bytes and any error encountered
func DecryptAES(encrypted_bytes, public_key *[]byte, padding_value byte) (decrypted *[]byte, err error) {
	// Get AES key from environment and return err if it was not present
	private_key := []byte(os.Getenv("SWU_CRYPTO_AES_PRIVATE_KEY"))
	if len(private_key) == 0 {
		return nil, errors.New("Missing AES private key")
	}

	cipher_block, err := aes.NewCipher(private_key)
	if err != nil {
		return nil, err
	}

	// Decrypt!
	mode := cipher.NewCBCDecrypter(cipher_block, *public_key)
	decrypted_bytes := make([]byte, len(*encrypted_bytes))
	mode.CryptBlocks(decrypted_bytes, *encrypted_bytes)
	// Trim the trim_byte byte from the decrypted value
	trimmed := bytes.Trim(decrypted_bytes, string(padding_value))
	// Return the byte pointer
	return &trimmed, err
}
