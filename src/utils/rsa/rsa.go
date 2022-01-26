package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func DecryptPKCS1v15(keyBlock []byte,ciphertext []byte)([]byte,error){
	pri, err := x509.ParsePKCS8PrivateKey(keyBlock)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, pri.(*rsa.PrivateKey), ciphertext)
}

func EncryptPKCS1v15(keyBlock []byte, plaintext []byte)([]byte, error)  {

	pub, err := x509.ParsePKCS1PublicKey(keyBlock)
	if nil != err{
		return nil, err
	}

	if len(plaintext) > pub.Size() - 11{
		return nil, fmt.Errorf("message too long for RSA public key size, must be less than %d", pub.Size() - 11)
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pub, plaintext)
}


func EncodePrivateKey(key *rsa.PrivateKey)string{
	pkcs8, _ := x509.MarshalPKCS8PrivateKey(key)
	privatePem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pkcs8,
		},
	)
	return string(privatePem)
}

func EncodePublicKey(key *rsa.PublicKey)string{
	pkcs1 := x509.MarshalPKCS1PublicKey(key)
	privatePem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pkcs1,
		},
	)
	return string(privatePem)
}