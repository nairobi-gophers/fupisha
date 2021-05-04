package store

import (
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents an authenticated user.
type User struct {
	ID                   uuid.UUID  `db:"id,omitempty"`
	Email                string     `db:"email,omitempty"`
	Password             string     `db:"password"`
	APIKey               *uuid.UUID `db:"api_key,omitempty"`
	ResetPasswordExpires *time.Time `db:"reset_password_expires,omitempty"`
	ResetPasswordToken   *uuid.UUID `db:"reset_password_token,omitempty"`
	VerificationExpires  time.Time  `db:"verification_expires"`
	VerificationToken    uuid.UUID  `db:"verification_token,omitempty"`
	Verified             bool       `db:"verified,omitempty"`
	CreatedAt            time.Time  `db:"created_at,omitempty"`
	UpdatedAt            time.Time  `db:"updated_at,omitempty"`
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
