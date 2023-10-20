package keycloak

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyTokenMiddleware verifies OIDC tokens
func (kc *Keycloak) VerifyTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawAccessToken, err := getAccessTokenFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		// Verify token with Keycloak
		resp, err := http.PostForm(kc.BaseUrl()+VerifyTokenPath,
			url.Values{
				"token":         {rawAccessToken},
				"client_id":     {kc.Client.ID},
				"client_secret": {kc.Client.Secret},
			},
		)

		if err != nil || resp.StatusCode != 200 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateTokenMiddleware checks for valid OIDC tokens with preloaded public key
func (kc *Keycloak) ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawAccessToken, err := getAccessTokenFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		// Parse and verify JWT token
		token, err := jwt.Parse(rawAccessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return kc.pk, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getAccessTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) < 2 {
		return "", fmt.Errorf("malformed token")
	}

	return parts[1], nil
}
