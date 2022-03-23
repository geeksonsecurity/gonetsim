package web

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"gonetsim/internal/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"software.sslmate.com/src/go-pkcs12"
	"time"
)

type Server struct {
	HttpPort    int
	HttpsPort   int
	httpServer  *http.Server
	httpsServer *http.Server
}

func (s *Server) Start(callback func(uri string, raw string)) {
	http.DefaultServeMux = new(http.ServeMux)
	s.startHttpServer(callback)
	s.startHttpsServer()
}

func (s *Server) Stop() {
	log.Printf("Stopping web servers\n")
	ctxShutDown, _ := context.WithTimeout(context.Background(), 2*time.Second)
	if s.httpServer != nil {
		log.Printf("Stopping HTTP server\n")
		s.httpServer.Shutdown(ctxShutDown)
	}
	if s.httpsServer != nil {
		log.Printf("Stopping HTTPs server\n")
		s.httpsServer.Shutdown(ctxShutDown)
	}
}

func (s *Server) startHttpServer(callback func(uri string, raw string)) {
	log.Printf("Starting HTTP server at port %d\n", s.HttpPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "@_@\n")

		schema := r.URL.Scheme
		requestURL := r.URL.String()
		if schema == "" {
			schema = "http"
			if r.TLS != nil {
				schema = "https"
			}
			requestURL = schema + "://" + r.Host + r.URL.Path
		}
		requestURL = r.Method + " " + requestURL + " " + r.Proto

		raw := requestURL + "\r\n"
		headers := ""
		// Loop over header names
		for name, values := range r.Header {
			// Loop over all values for the name.
			for _, value := range values {
				headers += fmt.Sprintf("%s: %s\r\n", name, value)
			}
		}
		raw += headers
		bodyB, _ := ioutil.ReadAll(r.Body)
		bodyStr := string(bytes.Replace(bodyB, []byte("\r"), []byte("\r\n"), -1))
		if bodyStr != "" {
			raw += "\r\n\r\n" + bodyStr
		}

		callback(requestURL, raw)
	})

	s.httpServer = &http.Server{
		Addr: fmt.Sprint(":", s.HttpPort),
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start http server %s\n", err.Error())
		}
	}()
}

func (s *Server) startHttpsServer() {
	log.Printf("Starting HTTPs server at port %d\n", s.HttpsPort)

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
		utils.WriteBinaryFile(p12Path, pfxBytes)
		utils.WriteBinaryFile(certPath, caCertificate.Raw)
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

	s.httpsServer = &http.Server{
		Addr:      fmt.Sprint(":", s.HttpsPort),
		TLSConfig: serverTLSConf,
	}

	go func() {
		defer s.httpsServer.Close()
		if err := s.httpsServer.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTPS server %s\n", err.Error())
		}
	}()
}
