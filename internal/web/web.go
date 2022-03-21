package web

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
)

func StartHttpServer(port int) {

	log.Printf("Starting HTTP server at port %d\n", port)

	http.HandleFunc("/", handleRequest())

	go func() {
		if err := http.ListenAndServe(fmt.Sprint(":", port), nil); err != nil {
			log.Fatalf("Failed to start http server %s\n", err.Error())
		}
	}()
}

func handleRequest() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "@_@\n")
		log.Printf("Responded to request %s %s%s", r.Method, r.Host, r.RequestURI)
	}
}

func StartHttpsServer(port int) {
	log.Printf("Starting HTTPs server at port %d\n", port)
	caCertificate, caBytes, caKey := GenerateCACertificate()

	caCert, _ := x509.ParseCertificate(caBytes)
	publicKeyDer, _ := x509.MarshalPKIXPublicKey(caCert.PublicKey)
	publicKeyBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDer,
	}
	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))
	log.Printf("CA Certificate (DER format):\n%s", publicKeyPem)

	// lets generate a key here and reuse it everytime!
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("Failed to generate caCertificate caKey %s\n", err.Error())
	}

	serverTLSConf := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert := SignCertificate(caCertificate, caKey, certPrivKey, info.ServerName)
			return &cert, nil
		},
	}

	s := &http.Server{
		Addr:      fmt.Sprint(":", port),
		TLSConfig: serverTLSConf,
	}

	go func() {
		defer s.Close()
		if err := s.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("Failed to start HTTPS server %s\n", err.Error())
		}
	}()
}
