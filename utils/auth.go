package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
)

func init() {
	PublicKey, PrivateKey = initKeys()
}

func initKeys() (*rsa.PublicKey, *rsa.PrivateKey) {
	// load the environment containing the public key
	err := godotenv.Load("../.env")
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

	certBytes, err := io.ReadAll(certFile)
	n := len(certBytes)
	if err != nil {
		if n == 0 {
			log.Println("certificate file not read")
		} else {
			log.Println("couldn't read the certificate file -> " + err.Error())
		}
		return nil, nil
	}
	certPem, _ := pem.Decode(certBytes) // decode the certificate

	privBytes, err := io.ReadAll(privFile)
	n = len(privBytes)
	if err != nil || n == 0 {
		if n == 0 {
			log.Println("certificate file not read")
		} else {
			log.Println("couldn't read private key file -> " + err.Error())
		}
		return nil, nil
	}
	privPem, _ := pem.Decode(privBytes) // decode the private key

	// extract public key from the certificate
	cert, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		log.Println("couldn't parse the certificate -> " + err.Error())
		return nil, nil
	}
	// extract private key
	key, err := x509.ParsePKCS8PrivateKey(privPem.Bytes)
	privKey := key.(*rsa.PrivateKey)
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
		if pubKey == nil {
			log.Println("nil public key")
		}
		if privKey == nil {
			log.Println("nil private key")
		}
		//return pubKey, privKey
		return pubKey, privKey
	}
}

/* returns nil if the user specified in the claims can access the target resource */
func Authorize(head *http.Header, target string) error {
	// extract token from header
	rawToken, err := ParseTokenFromHeader(head)
	if err != nil {
		return err
	}

	// validate the token
	claims, err := ValidToken(rawToken)
	if err != nil {
		return err
	}

	// get the username from the token
	if claims.Username != target {
		return errors.New("not authorized for that resource")
	}
	return nil
}
