package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func GenerateECKey() (key *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate ECDSA key: %s\n", err)
	}

	keyDer, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.Fatalf("Failed to serialize ECDSA key: %s\n", err)
	}

	keyBlock := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyDer,
	}

	keyFile, err := os.Create("ec_key.pem")
	if err != nil {
		log.Fatalf("Failed to open ec_key.pem for writing: %s", err)
	}
	defer func() {
		keyFile.Close()
	}()

	pem.Encode(keyFile, &keyBlock)

	return
}

func GenerateCert(pub, priv interface{}, cn string, ku x509.KeyUsage, filename string) {

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{"Sirius Service"},
		},
		KeyUsage:  ku,
		NotBefore: time.Now().Add(-time.Hour * 24 * 365),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),
	}
	certDer, err := x509.CreateCertificate(
		rand.Reader, &template, &template, pub, priv,
	)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s\n", err)
	}

	certBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDer,
	}

	certFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to open '%s' for writing: %s", filename, err)
	}
	defer func() {
		certFile.Close()
	}()

	pem.Encode(certFile, &certBlock)
}
