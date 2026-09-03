package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// PKCS7Padding PKCS7填充
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// PKCS7Unpadding PKCS7去除填充
func PKCS7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("invalid padding")
	}
	return data[:length-padding], nil
}

// Encrypt AES加密 (CBC模式)
func Encrypt(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	plainText = PKCS7Padding(plainText, blockSize)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt AES解密 (CBC模式)
func Decrypt(cipherText string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(data) < blockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := data[:blockSize]
	data = data[blockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	return PKCS7Unpadding(data)
}

// EncryptGCM AES加密 (GCM模式, 更安全)
func EncryptGCM(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptGCM AES解密 (GCM模式)
func DecryptGCM(cipherText string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, cipherTextData := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, cipherTextData, nil)
}

// ValidateKey 验证AES密钥长度
func ValidateKey(key []byte) error {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return errors.New("AES key length must be 16, 24, or 32 bytes")
	}
	return nil
}