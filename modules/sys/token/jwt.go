package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/iagapie/go-spring/modules/sys/config"
	"io/ioutil"
	"log"
	"time"
)

type jwtToken struct {
	privateKey []byte
	publicKey  []byte
}

type JWTOption interface {
	apply(jt *jwtToken)
}

type jwtOption func(jt *jwtToken)

func (fn jwtOption) apply(jt *jwtToken) {
	fn(jt)
}

func WithJWTKeys(jwtKeys config.JWTKeys) JWTOption {
	return jwtOption(func(jt *jwtToken) {
		WithPrivateKeyFile(jwtKeys.Private).apply(jt)
		WithPublicKeyFile(jwtKeys.Public).apply(jt)
	})
}

func WithPrivateKey(privateKey []byte) JWTOption {
	return jwtOption(func(jt *jwtToken) {
		jt.privateKey = privateKey
	})
}

func WithPublicKey(publicKey []byte) JWTOption {
	return jwtOption(func(jt *jwtToken) {
		jt.publicKey = publicKey
	})
}

func WithPrivateKeyFile(privateKeyFile string) JWTOption {
	prvKey, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		log.Fatalln(err)
	}
	return WithPrivateKey(prvKey)
}

func WithPublicKeyFile(publicKeyFile string) JWTOption {
	pubKey, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		log.Fatalln(err)
	}
	return WithPublicKey(pubKey)
}

func New(opts ...JWTOption) Token {
	jt := new(jwtToken)
	for _, opt := range opts {
		opt.apply(jt)
	}
	return jt
}

func (jt *jwtToken) Create(ttl time.Duration, content interface{}) (string, error) {
	if jt.privateKey == nil {
		return "", errors.New("create: private key is nil")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(jt.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content             // Our custom data.
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["nbf"] = now.Unix()          // The time before which the token must be disregarded.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (jt *jwtToken) Validate(token string) (interface{}, error) {
	if jt.publicKey == nil {
		return nil, errors.New("validate: public key is nil")
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(jt.publicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("validate: invalid")
	}

	return claims["dat"], nil
}
