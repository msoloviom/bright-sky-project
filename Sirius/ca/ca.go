package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

type ECDSASignature struct {
	R *big.Int
	S *big.Int
}

func GenerateECKey(fn string) (key *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
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

	keyFile, err := os.Create(fn)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %s", fn, err)
	}
	defer func() {
		keyFile.Close()
	}()

	pem.Encode(keyFile, &keyBlock)

	return
}

func GenerateCert(pub, priv interface{}, cert_signer *x509.Certificate, cn string, ku x509.KeyUsage, filename string) {
	sn, _ := ioutil.ReadFile("serial")
	s, _ := strconv.Atoi(string(sn))
	template := x509.Certificate{
		SerialNumber: big.NewInt(int64(s)),
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{"Sirius Service"},
		},
		KeyUsage:  ku,
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),
	}
	certDer, err := x509.CreateCertificate(
		rand.Reader, &template, cert_signer, pub, priv,
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
	ioutil.WriteFile("serial", []byte(strconv.Itoa(s+1)), 0644)
}

func VerifySignature(b64signature, pemcert string, data []byte) bool {
	derSignature, err := base64.StdEncoding.DecodeString(b64signature)
	if err != nil {
		return false
	}
	sig := ECDSASignature{}
	_, err = asn1.Unmarshal(derSignature, &sig)
	if err != nil {
		fmt.Print(err)
		return false
	}
	hash := sha512.Sum384(data)
	certBlock, rest := pem.Decode([]byte(pemcert))
	if len(rest) > 0 {
		return false
	}
	certObj, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return false
	}
	pubKey := certObj.PublicKey.(*ecdsa.PublicKey)
	return ecdsa.Verify(pubKey, hash[:], sig.R, sig.S)
}

func main() {
	switch os.Args[1] {
	case "-g":
		if len(os.Args) < 3 {
			log.Fatal("No CN provided!")
		}
		fn := os.Args[2]
		log.Printf("Generating an ECDSA P-384 Private Key to %s.key", fn)
		keyPem, _ := ioutil.ReadFile("sirius.key")
		certPem, _ := ioutil.ReadFile("sirius.crt")
		keyBlock, _ := pem.Decode(keyPem)
		certBlock, _ := pem.Decode(certPem)

		priv, _ := x509.ParseECPrivateKey(keyBlock.Bytes)
		cert, _ := x509.ParseCertificate(certBlock.Bytes)

		ECKey := GenerateECKey(fn + ".key")
		GenerateCert(&ECKey.PublicKey, priv, cert, fn, x509.KeyUsageDigitalSignature, fn+".crt")
	case "-s":
		var key, data string
		if len(os.Args) < 4 {
			log.Fatal("Invalid params!")
		}
		key = os.Args[2]
		data = os.Args[3]
		//cert := "Han Solo.crt"
		//certPem, _ := ioutil.ReadFile(cert)

		log.Printf("Signing %s with key %s", data, key)
		keyPem, err := ioutil.ReadFile(key)
		if err != nil {
			log.Fatal(err)
		}
		dataRaw, err := ioutil.ReadFile(data)
		if err != nil {
			log.Fatal(err)
		}
		keyBlock, _ := pem.Decode(keyPem)
		priv, err := x509.ParseECPrivateKey(keyBlock.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		hash := sha512.Sum384(dataRaw)
		r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
		if err != nil {
			log.Fatal(err)
		}
		sig := ECDSASignature{r, s}
		sigDer, err := asn1.Marshal(sig)
		if err != nil {
			log.Fatal(err)
		}
		sigb64 := base64.StdEncoding.EncodeToString(sigDer)

		fmt.Println(sigb64)
		//fmt.Print(VerifySignature(string(sigb64), string(certPem), dataRaw))
	case "-v":
		certFile := os.Args[2]
		//fmt.Println(certFile)
		signatureFile := os.Args[3]
		dataFile := os.Args[4]
		//sig := ECDSASignature{}
		//keyPem, err := ioutil.ReadFile("Han Solo.key")
		certPem, _ := ioutil.ReadFile(certFile)
		data, _ := ioutil.ReadFile(dataFile)
		signature, _ := ioutil.ReadFile(signatureFile)
		//fmt.Println(certPem)
		//fmt.Println(data)
		//fmt.Println(signature)
		fmt.Print(VerifySignature(string(signature), string(certPem), data))

	}
}
