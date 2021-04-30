package encoding

import (
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"log"

	"github.com/gofrs/uuid"
	nanoid "github.com/matoous/go-nanoid"
)

//Encode encodes the uuid to a base64 string that is url-safe.
func Encode(id uuid.UUID) string {
	return b64.RawURLEncoding.EncodeToString(id.Bytes())
}

//Decode decodes a base64 string to a raw uuid.
func Decode(id string) (uuid.UUID, error) {
	dec, err := b64.RawURLEncoding.DecodeString(id)

	if err != nil {
		return uuid.UUID{}, err
	}

	decoded, err := uuid.FromBytes(dec)
	if err != nil {
		return uuid.UUID{}, err
	}

	return decoded, nil
}

//GenUniqueParam returns a random param but unique key.
func GenUniqueParam(alphanumeric string, len int) (string, error) {
	return nanoid.Generate(alphanumeric, len)
}

//GenUniqueID returns a random but unique id.
func GenUniqueID() uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate id :%s", err)
	}
	return id
}

//GenHexKey generates a crypto-random key with byte length len and hex-encodes it to a string.
func GenHexKey(len int) string {
	bytes := make([]byte, len)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(bytes)
}
