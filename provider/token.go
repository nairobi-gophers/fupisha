package provider

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nairobi-gophers/fupisha/config"
)

//JWTService defines an implementation of our JWT authentication service.
type JWTService interface {
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

//NewJWTService configures and returns a JWT authentication instance.
func NewJWTService(secret string) (JWTService, error) {

	if len([]byte(secret)) < 32 {
		return nil, errors.New("jwt: secret too short")
	}

	return &service{secret: []byte(secret)}, nil
}

// Encode a claim into a JWT
func (s *service) Encode(uid string) (string, error) {

	// Create the Claims
	c := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.cfg.JWT.ExpireDelta)).Unix(),
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
