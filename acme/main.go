package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	domain := flag.String("domain", "", "Obtain certificate for this domain")
	flag.Parse()

	if "" == *domain {
		log.Fatal("You did not set value for domain")
	}

	myUser := createUser()
	caCertPool := createCertPool()

	// This CA URL is configured for a local dev instance of Pebble running in Docker
	config := lego.NewConfig(myUser)
	config.CADirURL = "https://127.0.0.1:14000/dir"
	config.Certificate.KeyType = certcrypto.RSA2048
	config.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            caCertPool,
		},
	}
	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// We specify an HTTP port of 5002 and an TLS port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that you still
	// need to proxy challenge traffic to port 5002 and 5001.
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "5002"))
	if err != nil {
		log.Fatal(err)
	}
	//err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "5001"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{*domain},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("cert.pem", certificates.Certificate, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("key.pem", certificates.PrivateKey, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func createUser() *MyUser {
	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	return &MyUser{
		Email: "you@yours.com",
		key:   privateKey,
	}
}

func createCertPool() *x509.CertPool {
	f, err := os.Open("./certs/pebble.minica.pem")
	if err != nil {
		log.Fatal(err)
	}

	cert, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	return caCertPool
}
