package aws

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

var key, _ = hex.DecodeString("93d9f6846b629edb2bdc4466af627d998496cb0c08f9cf043de68d6b25aa9693")

func Encrypt(plaintext []byte) (string, error) {
	iv := make([]byte, 12)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	ciphertext := aesgcm.Seal(nil, iv, plaintext, nil)
	tag := ciphertext[len(ciphertext)-16:]
	ct := ciphertext[:len(ciphertext)-16]
	
	ivB64 := base64.StdEncoding.EncodeToString(iv)
	return fmt.Sprintf("%s::%s::%s", ivB64, hex.EncodeToString(tag), hex.EncodeToString(ct)), nil
}

func Decrypt(encrypted string) ([]byte, error) {
	parts := strings.Split(encrypted, "::")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format")
	}
	
	iv, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	
	tag, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	
	ct, err := hex.DecodeString(parts[2])
	if err != nil {
		return nil, err
	}
	
	ciphertext := append(ct, tag...)
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	return aesgcm.Open(nil, iv, ciphertext, nil)
}
