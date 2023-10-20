package keycloak

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
)

type Keycloak struct {
	Conf   *Conf
	Client *Client
	pk     *rsa.PublicKey
}

type Conf struct {
	Host   string
	Port   string
	UseSSl bool
	Realm  string
}

type Client struct {
	ID     string
	Secret string
}

type Jwks struct {
	Keys []struct {
		Kty string `json:"kty"`
		E   string `json:"e"`
		N   string `json:"n"`
	} `json:"keys"`
}

const (
	CertPath        = "/protocol/openid-connect/certs"
	AuthPath        = "/protocol/openid-connect/auth"
	TokenPath       = "/protocol/openid-connect/token"
	VerifyTokenPath = "/protocol/openid-connect/token/introspect"
)

func (kc *Keycloak) LoadPublicKey() error {
	pk, err := kc.publicKey()
	if err != nil {
		return err
	}

	kc.pk = pk
	return nil
}

func (kc *Keycloak) publicKey() (*rsa.PublicKey, error) {
	// Fetch public key from Keycloak
	resp, err := http.Get(kc.BaseUrl() + CertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWK: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWK: %s", err)
	}

	var jwks Jwks
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWK: %s", err)
	}

	// Assuming the first key is the one we want
	eBytes, err := decodeSegment(jwks.Keys[0].E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWK exponent: %s", err)
	}

	nBytes, err := decodeSegment(jwks.Keys[0].N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWK modulus: %s", err)
	}

	return &rsa.PublicKey{
		E: int(new(big.Int).SetBytes(eBytes).Int64()), // Convert the RSA exponent from big.Int to int; typically small, often 65537 (<= 2^16)
		N: new(big.Int).SetBytes(nBytes),
	}, nil
}

func decodeSegment(seg string) ([]byte, error) {
	padLen := 4 - (len(seg) % 4)
	if padLen < 4 {
		seg += strings.Repeat("=", padLen)
	}
	return base64.URLEncoding.DecodeString(seg)
}

func (kc *Keycloak) BaseUrl() string {
	if kc.Conf.UseSSl {
		return "https://" + kc.Conf.Host + ":" + kc.Conf.Port + "/realms/" + kc.Conf.Realm
	}
	return "http://" + kc.Conf.Host + ":" + kc.Conf.Port + "/realms/" + kc.Conf.Realm
}
