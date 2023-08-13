package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := genKey(); err != nil {
		log.Fatalln(err)
	}
}

func genKey() error {
	// Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	// Create a file for the private key information in PEM format
	privateKeyFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("creating private key file: %w", err)
	}
	defer privateKeyFile.Close()

	// Construct a PEM block to represent the private key.
	privateBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Write the private key to the private key file.
	if err := pem.Encode(privateKeyFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding private key: %w", err)
	}

	// ----------------------------------------------------------

	// Create a file for the public key information in PEM format
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("creating public key file: %w", err)
	}
	defer publicFile.Close()

	// Marshal the public key from the private key to PKIX format.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	// Construct a PEM block to represent the public key.
	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// Write the public key to the public key file.
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding public key: %w", err)
	}

	fmt.Println("Keys generated successfully.")

	return nil
}
