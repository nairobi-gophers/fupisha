package encoding

import (
	b64 "encoding/base64"

	"github.com/gofrs/uuid"
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
