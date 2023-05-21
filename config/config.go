package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/dundunlabs/omniauth"
	"github.com/uptrace/bun"
)

type Config struct {
	DB    *bun.DB
	OAuth map[string]omniauth.OmniAuth
	JWT   *JWTConfig
}

type JWTConfig struct {
	PrivateKeyPEM string
	PublicKeyPEM  string
}

func (jwt *JWTConfig) PrivateKey() (*rsa.PrivateKey, error) {
	privBlock, _ := pem.Decode([]byte(jwt.PrivateKeyPEM))
	return x509.ParsePKCS1PrivateKey(privBlock.Bytes)
}

func (jwt *JWTConfig) PublicKey() (*rsa.PublicKey, error) {
	pubBlock, _ := pem.Decode([]byte(jwt.PublicKeyPEM))
	return x509.ParsePKCS1PublicKey(pubBlock.Bytes)
}
