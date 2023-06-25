package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	// Generate a new private key.
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	// Create a file for the private key information in PEM form
	// privFile, err := os.Create("private.pem")
	// if err != nil {
	// 	return fmt.Errorf("creating private file: %w", err)
	// }
	// defer privFile.Close()

	// Construct a PEM block for the private key.
	// privBlock := pem.Block{
	// 	Type:  "PRIVATE KEY",
	// 	Bytes: x509.MarshalPKCS1PrivateKey(privkey),
	// }

	// Write the private key to the private key file.
	// if err := pem.Encode(privFile, &privBlock); err != nil {
	// 	return fmt.Errorf("encoding to private file: %w", err)
	// }

	// Create a file for the public key information in PEM form.
	// pubFile, err := os.Create("public.pem")
	// if err != nil {
	// 	return fmt.Errorf("creating public file: %w", err)
	// }
	// defer pubFile.Close()

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privkey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	// Construct a PEM block for the public key.
	pubBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	fmt.Print("=====================================\n\n")

	// Write the public key to the public key file.
	if err := pem.Encode(os.Stdout, &pubBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	fmt.Print("\n=====================================\n\n")

	// =============================

	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "123456789",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	token.Header["kid"] = "kid1"

	str, err := token.SignedString(privkey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println(str)

	// fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n", token)
	return nil
}
