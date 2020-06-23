package model

import (
	"time"

	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User represents an authenticated user.
type User struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	Name                 string             `bson:"name,omitempty"`
	CreatedAt            time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt            time.Time          `bson:"updatedAt, omitempty"`
	Email                string             `bson:"email,omitempty"`
	Password             string             `bson:"password,omitempty"`
	APIKey               uuid.UUID          `bson:"apiKey,omitempty"`
	ResetPasswordExpires time.Time          `bson:"resetPasswordExpires,omitempty"`
	ResetPasswordToken   string             `bson:"resetPasswordToken,omitempty"`
	VerificationExpires  time.Time          `bson:"verificationExpires,omitempty"`
	VerificationToken    uuid.UUID          `bson:"verificationToken,omitempty"`
	Verified             bool               `bson:"verified,omitempty"`
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
