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
	//if err := genKey(); err != nil {
	if err := genToken(); err != nil {
		log.Fatalln(err)
	}
}

func genToken() error {

	// Generating a token requires defining a set of claims. In this application's case,
	// we only care about defining the subject and the user in question and the roles
	// they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the jwt
	// sub (subject): Subject of the jwt (the user)
	// aud (audience): Recipient for which the jwt is intended
	// exp (expiration time): Time after which the jwt expires
	// nbf (not before time): Time before which the jwt must not be accepted for processing
	// iat (issued at time): Time at which the jwt was issued; can be used to determine age of the jwt
	// jti (jwt id): Unique identifier; can be used to prevent the jwt from being replayed (allows a token to be used only once)
	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "12345678",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(jwt.SigningMethodRS256.Name), claims)
	token.Header["kid"] = "kid1"

	// Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	str, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("************* TOKEN *************")
	fmt.Println(str)
	fmt.Println()

	// ----------------------------------------------------------

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	if err := pem.Encode(os.Stdout, &publicBlock); err != nil {
		return fmt.Errorf("encoding public key: %w", err)
	}

	fmt.Println()

	// ----------------------------------------------------------
	// Validate JWT with Public Key

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

	var clm struct {
		jwt.RegisteredClaims
		Roles []string
	}

	kf := func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	}

	tkn, err := parser.ParseWithClaims(str, &clm, kf)
	if err != nil {
		return fmt.Errorf("parsing with claims: %w", err)
	}

	if !tkn.Valid {
		return fmt.Errorf("token is invalid")
	}

	fmt.Println("TOKEN VALIDATED SUCCESSFULLY")
	fmt.Printf("%#v\n", clm)

	return nil
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
