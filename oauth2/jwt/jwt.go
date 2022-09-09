package jwt

import (
	"barton.top/btgo/pkg/oauth2"
	"barton.top/btgo/pkg/oauth2/jws"
	"crypto/rsa"
	"fmt"
	goauth2 "golang.org/x/oauth2"
	"time"
)

// Config is the configuration for using JWT to fetch tokens,
// commonly known as "two-legged OAuth 2.0".
type Config struct {
	Platform string

	// Username is the OAuth client identifier used when communicating with
	// the configured OAuth provider.
	Username string

	// PrivateKey contains the contents of an RSA private key or the
	// contents of a PEM file that contains a private key. The provided
	// private key is used to sign JWT payloads.
	// PEM containers with a passphrase are not supported.
	// Use the following command to convert a PKCS 12 file into a PEM.
	//
	//    $ openssl pkcs12 -in key.p12 -out key.pem -nodes
	//
	PrivateKey []byte

	// PrivateKeyID contains an optional hint indicating which key is being
	// used.
	PrivateKeyID string

	// Subject is the optional user to impersonate.
	Subject string

	// Scopes optionally specifies a list of requested permission scopes.
	Scopes []string

	// Expires optionally specifies how long the token is valid for.
	Expires time.Duration
}

// TokenSource returns a JWT TokenSource using the configuration
// in c and the HTTP client from the provided context.
func (c *Config) TokenSource() goauth2.TokenSource {
	return goauth2.ReuseTokenSource(nil, &jwtSource{conf: c})
}

type jwtSource struct {
	conf *Config
}

func (ts *jwtSource) Token() (*oauth2.Token, error) {
	pk, err := oauth2.ParseKey[rsa.PrivateKey](ts.conf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not parse key: %v", err)
	}

	// create access token
	iat := time.Now()
	exp := iat.Add(ts.conf.Expires)
	cs := &jws.ClaimSet{
		Iss: ts.conf.Platform,
		Sub: ts.conf.Username,
		Aud: "",
		Iat: iat.Unix(),
		Exp: exp.Unix(),
	}
	hdr := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}
	accessToken, err := jws.Encode(hdr, cs, pk)
	if err != nil {
		return nil, fmt.Errorf("could not encode JWT: %v", err)
	}

	// create refresh token
	cs = &jws.ClaimSet{
		Iss: ts.conf.Username,
		Sub: ts.conf.Username,
		Aud: "",
		Iat: iat.Unix(),
		Exp: exp.Unix() + (int64)(time.Hour*24*365*5),
	}
	refreshToken, err := jws.Encode(hdr, cs, pk)
	if err != nil {
		return nil, fmt.Errorf("could not encode JWT: %v", err)
	}

	return &oauth2.Token{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "Bearer", Expiry: exp}, nil
}
