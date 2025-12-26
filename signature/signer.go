package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func SignPayload(payload string, privateKey []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("invalid private key PEM")
	}

	var rsaPrivateKey *rsa.PrivateKey
	var err error

	if key, err1 := x509.ParsePKCS1PrivateKey(block.Bytes); err1 == nil {
		rsaPrivateKey = key
	} else if key, err2 := x509.ParsePKCS8PrivateKey(block.Bytes); err2 == nil {
		rsaPrivateKey = key.(*rsa.PrivateKey)
	} else {
		return "", err
	}

	hashed := sha256.Sum256([]byte(payload))

	signature, err := rsa.SignPKCS1v15(
		rand.Reader,
		rsaPrivateKey,
		crypto.SHA256,
		hashed[:],
	)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
