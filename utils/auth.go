package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	SignMethod = jwt.SigningMethodRSA{
		Name: "RSA-2048",
		Hash: crypto.SHA256,
	}

	PublicKey, PrivateKey = initKeys()
)

func initKeys() (*rsa.PublicKey, *rsa.PrivateKey) {
	// load the environment containing the public key
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("couldn't load environment file -> " + err.Error())
		return nil, nil
	}

	certPath := os.Getenv("CERT")
	privPath := os.Getenv("KEY")
	certFile, err := os.Open(certPath)
	if err != nil {
		log.Println("couldn't open certificate -> " + err.Error())
		return nil, nil
	}
	privFile, err := os.Open(privPath)
	if err != nil {
		log.Println("couldn't open private key file -> " + err.Error())
		return nil, nil
	}

	var certBytes []byte
	var privBytes []byte
	_, err = certFile.Read(certBytes)
	if err != nil {
		log.Println("couldn't read the certificate file -> " + err.Error())
		return nil, nil
	}
	_, err = privFile.Read(privBytes)
	if err != nil {
		log.Println("couldn't read private key file -> " + err.Error())
	}

	// extract public key from the certificate
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		log.Println("couldn't parse the certificate -> " + err.Error())
		return nil, nil
	}
	// extract private key
	privKey, err := x509.ParsePKCS1PrivateKey(privBytes)
	if err != nil {
		log.Println("couldn't parse the private key -> " + err.Error())
		return nil, nil
	}

	pubKey := cert.PublicKey.(*rsa.PublicKey)

	// validate the certificate key with the private key
	if pubKey.N.Cmp(privKey.N) != 0 {
		log.Println("Certificate key doesn't match private key")
		return nil, nil
	} else {
		return pubKey, privKey
	}
}

func VerifyToken(t *jwt.Token) (any, error) { return PublicKey, nil }
