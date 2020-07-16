package provider

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nairobi-gophers/fupisha/internal/config"
)

//Service defines an implementation of our JWT authentication service.
type Service interface {
	Encode(uid string) (token string, err error)
	Decode(token string) (userID string, issuedAt time.Time, err error)
}

type service struct {
	secret []byte
	cfg    config.Config
}

// Claims is our custom metadata, which will be hashed
// and sent as the second segment in our JWT.
type Claims struct {
	jwt.StandardClaims
	UserID string
}

//NewService configures and returns a JWT authentication instance.
func NewService(secret string) (Service, error) {
	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("jwt: failed to decode jwt secret from string %s: %w", secret, err)
	}

	if len(secretBytes) < 32 {
		return nil, errors.New("jwt: secret too short")
	}

	return &service{secret: secretBytes}, nil
}

// Encode a claim into a JWT
func (s *service) Encode(uid string) (string, error) {
	// Create the Claims
	c := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(1)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "fupisha",
		},
		UserID: uid,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// Sign token and return
	return token.SignedString(s.secret)
}

//Decode verifies the JWT string using the given secret key,
//on success it decodes it and returns the user ID string.
func (s *service) Decode(tkn string) (string, time.Time, error) {
	tokenType, err := jwt.ParseWithClaims(tkn, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt: unexpected signing method")
		}
		return s.secret, nil
	})

	var uid string

	if err != nil {
		return uid, time.Time{}, fmt.Errorf("jwt: ParseWithClaims failed: %w", err)
	}

	c, ok := tokenType.Claims.(*Claims)
	if !ok {
		return uid, time.Time{}, errors.New("jwt: failed to get token claims")
	}

	if c.UserID == "" {
		return uid, time.Time{}, errors.New("jwt: UserID claim is not valid")
	}

	if c.IssuedAt == 0 {
		return uid, time.Time{}, errors.New("jwt: IssuedAt claim is not valid")
	}

	if c.Issuer == "" || c.Issuer != "fupisha" {
		return uid, time.Time{}, errors.New("jwt: Issuer claim is not valid")
	}

	return c.UserID, time.Unix(c.IssuedAt, 0), nil
}
