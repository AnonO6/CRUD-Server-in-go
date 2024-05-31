package main

import (
	"crypto/rsa"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// loadRSAKeys loads the RSA keys from PEM files
func loadRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKeyBytes, err := os.ReadFile("private_key.pem")
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't load private_key.pem: %v", err)
	}
	
	privateKey,err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	publicKeyBytes, err := os.ReadFile("public_key.pem")
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't load public_key.pem: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	return privateKey, publicKey, nil
}


func generateJWT(userID int, privateKey *rsa.PrivateKey, expiresInSeconds int64) (string, error) {
	// This is for custom duration time, if set by client however cap is 24hrs
	expirationTime := time.Hour * 24
	if expiresInSeconds > 0 {
		expirationTime = time.Duration(expiresInSeconds) * time.Second
		if expirationTime > 24*time.Hour {
			expirationTime = 24 * time.Hour
		}
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
