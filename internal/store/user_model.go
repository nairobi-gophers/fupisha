package model

import (
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents an authenticated user.
type User struct {
	ID                   string    `bson:"_id,omitempty" db:"id,omitempty"`
	Name                 string    `bson:"name,omitempty" db:"name,omitempty"`
	CreatedAt            time.Time `bson:"createdAt,omitempty" db:"created_at,omitempty"`
	UpdatedAt            time.Time `bson:"updatedAt, omitempty" db:"updated_at,omitempty"`
	Email                string    `bson:"email,omitempty" db:"email,omitempty"`
	Password             string    `bson:"password,omitempty" db:"password"`
	APIKey               uuid.UUID `bson:"apiKey,omitempty" db:"api_key,omitempty"`
	ResetPasswordExpires time.Time `bson:"resetPasswordExpires,omitempty" db:"reset_password_expires,omitempty"`
	ResetPasswordToken   string    `bson:"resetPasswordToken,omitempty" db:"reset_password_token,omitempty"`
	VerificationExpires  time.Time `bson:"verificationExpires,omitempty" db:"verification_expires"`
	VerificationToken    uuid.UUID `bson:"verificationToken,omitempty" db:"verification_token,omitempty"`
	Verified             bool      `bson:"verified,omitempty" db:"verified,omitempty"`
}

//HashPassword hashes the user password using bcrypt hash function
func (u *User) HashPassword() error {

	pwd := []byte(u.Password)

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}

//Compare compares the password hash against the passed in password string
func (u User) Compare(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
