package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func RSAEncrypt(plain string, publicKey string) ([]byte, error) {
	// Tạo khóa RSA
	// privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	// if err != nil {
	// 	fmt.Println("Lỗi tạo khóa:", err)
	// 	return nil, nil
	// }

	// // Trích xuất khóa công khai từ khóa RSA
	// publicKeys := &privateKey.PublicKey

	// // Chuyển đổi khóa công khai thành định dạng PEM
	// publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKeys)
	// publicKeyPEM := pem.EncodeToMemory(&pem.Block{
	// 	Type:  "RSA PUBLIC KEY",
	// 	Bytes: publicKeyBytes,
	// })

	// // Hiển thị khóa công khai
	// publicKey = string(publicKeyPEM)

	msg := []byte(plain)

	// Decode public key
	pubBlock, _ := pem.Decode([]byte(publicKey))
	if pubBlock == nil {
		return nil, errors.New("failed to decode public key")
	}

	pubKey, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, err
	}

	encryptedMsg, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, msg)
	if err != nil {
		return nil, err
	}

	return encryptedMsg, nil
}
