package web

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"software.sslmate.com/src/go-pkcs12"
)

func StartWebServers(httpPort, httpsPort int, callback func(protocol string, uri string)) {
	startHttpServer(httpPort, callback)
	startHttpsServer(httpsPort)
}

func startHttpServer(httpPort int, callback func(protocol string, uri string)) {
	log.Printf("Starting HTTP server at port %d\n", httpPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "@_@\n")
		protocol := "HTTP"
		if r.TLS != nil {
			protocol = "HTTPS"
		}
		callback(protocol, fmt.Sprintf("%s %s%s", r.Method, r.Host, r.RequestURI))
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprint(":", httpPort), nil); err != nil {
			log.Fatalf("Failed to start http server %s\n", err.Error())
		}
	}()
}

func startHttpsServer(port int) {
	log.Printf("Starting HTTPs server at port %d\n", port)

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	var p12Path = filepath.Join(exPath, "ca.p12")
	var certPath = filepath.Join(exPath, "ca.crt")

	var caCertificate *x509.Certificate
	var caKey *rsa.PrivateKey

	file, err := os.Open(p12Path)
	if errors.Is(err, os.ErrNotExist) {
		caCertificate, caKey = GenerateCACertificate()
		pfxBytes, err := pkcs12.Encode(rand.Reader, caKey, caCertificate, []*x509.Certificate{}, pkcs12.DefaultPassword)
		if err != nil {
			log.Fatal(err)
		}
		writeFile(p12Path, pfxBytes)
		writeFile(certPath, caCertificate.Raw)
	} else {
		pfxBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		caKey2, caCertificate2, err := pkcs12.Decode(pfxBytes, pkcs12.DefaultPassword)
		if err != nil {
			log.Fatal(err)
		}
		caCertificate = caCertificate2
		caKey = caKey2.(*rsa.PrivateKey)
	}

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

func writeFile(path string, content []byte) {
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bytesWritten, err := file.Write(content)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes to %s.\n", bytesWritten, path)
}
