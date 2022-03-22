package web

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	mathrand "math/rand"
	"time"
)

func GenerateCACertificate() (caCert *x509.Certificate, caKey *rsa.PrivateKey) {
	// set up our CA certificate
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Gonetsim"},
			Country:       []string{"CH"},
			Province:      []string{""},
			Locality:      []string{"Zurich"},
			StreetAddress: []string{"Paradeplatz"},
			PostalCode:    []string{"8000"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("Failed to generate CA key %s\n", err.Error())
	}

	// create the CA
	caFinal, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		log.Fatalf("Failed to generate CA certificate %s\n", err.Error())
	}
	cert, err := x509.ParseCertificate(caFinal)
	return cert, caPrivKey
}

func SignCertificate(ca *x509.Certificate, caPrivKey *rsa.PrivateKey, certPrivKey *rsa.PrivateKey, commonName string) (serverTLSConf tls.Certificate) {

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(int64(mathrand.Uint64())),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		log.Fatalf("Failed to generate certificate %s\n", err.Error())
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		log.Fatalf("Failed to generate certificate keypair %s\n", err.Error())
	}

	return serverCert
}
