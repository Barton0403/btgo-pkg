package oauth2

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
	goauth2 "golang.org/x/oauth2"
	"net/http"
)

type Token = goauth2.Token
type ClientInfo = oauth2.ClientInfo

// ParseKey converts the binary contents of a private key file
// to an *T. It detects whether the private key is in a
// PEM container or not. If so, it extracts the the private key
// from PEM container before conversion. It only supports PEM
// containers with no passphrase.
func ParseKey[T rsa.PrivateKey | rsa.PublicKey](key []byte) (*T, error) {
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("private key should be a PEM or plain PKSC1 or PKCS8; parse error: %v", err)
		}
	}
	parsed, ok := parsedKey.(*T)
	if !ok {
		return nil, errors.New("private key is invalid")
	}
	return parsed, nil
}

type Server interface {
	SetClientStorage(store oauth2.ClientStore)
	SetTokenStorage(store oauth2.TokenStore)
	SetPasswordAuthorizationHandler(handler server.PasswordAuthorizationHandler)
	SetUserAuthorizationHandler(handler server.UserAuthorizationHandler)
	ValidationBearerToken(r *http.Request) (oauth2.TokenInfo, error)
	HandleTokenRequest(w http.ResponseWriter, r *http.Request) error
	HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error
	GetManager() oauth2.Manager
	Logout(ctx context.Context, id string)
}
