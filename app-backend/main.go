package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/neox5/auth-proto/app-backend/pkg/keycloak"
	"golang.org/x/oauth2"
)

var (
	clientID     = "auth-proto-app"
	clientSecret = "ZYRYbJfeVsjElaNthoku1dd6eh4s5tDF"
)
var oauth2Config = oauth2.Config{
	ClientID:     clientID,     // Your client ID
	ClientSecret: clientSecret, // Your client secret
	RedirectURL:  "http://localhost:3000/auth/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "http://localhost:8080/realms/keycloak-test/protocol/openid-connect/auth",
		TokenURL: "http://localhost:8080/realms/keycloak-test/protocol/openid-connect/token",
	},
	Scopes: []string{"openid", "profile", "email"},
}

var keycloakConf = keycloak.Conf{
	Host:   "localhost",
	Port:   "8080",
	UseSSl: false,
	Realm:  "keycloak-test",
}

var verifier *oidc.IDTokenVerifier

var kc *keycloak.Keycloak

func init() {
	kc = &keycloak.Keycloak{
		Conf: &keycloak.Conf{
			Host:   "localhost",
			Port:   "8080",
			UseSSl: false,
			Realm:  "keycloak-test",
		},
		Client: &keycloak.Client{
			ID:     clientID,
			Secret: clientSecret,
		},
	}

	// Fetch public key from Keycloak
	err := kc.LoadPublicKey()
	if err != nil {
		log.Fatalf("Keycloak: Could not load public key: %v", err)
	}
}

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, kc.BaseUrl())
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	oidcConfig := &oidc.Config{ClientID: clientID}
	verifier = provider.Verifier(oidcConfig)

	r := chi.NewRouter()

	r.Get("/public", PublicHandler)
	r.With(kc.VerifyTokenMiddleware).Get("/private", PrivateHandler)
	r.With(kc.ValidateTokenMiddleware).Get("/private2", PrivateHandler2)
	r.Get("/auth", AuthHandler)
	r.Get("/auth/callback", AuthCallbackHandler)

	http.ListenAndServe(":3000", r)
}

// PublicHandler serves public content
func PublicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Public content"))
}

// PrivateHandler serves private content
func PrivateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Private content"))
}

// PrivateHandler serves private content
func PrivateHandler2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a private2 route"))
}

// AuthHandler redirects to Keycloak login
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauth2Config.AuthCodeURL("state"), http.StatusFound)
}

// AuthCallbackHandler handles the callback from Keycloak
func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}
	_, err = verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Display tokens to the user
	resp := fmt.Sprintf("ID Token: %s<br/>Access Token: %s", rawIDToken, oauth2Token.AccessToken)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(resp))
}
