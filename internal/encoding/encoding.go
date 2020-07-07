package encoding

import "github.com/gofrs/uuid"

//Encode encodes the uuid to a base64 string that is url-safe.
func Encode(id uuid.UUID) string {
	return b64.RawURLEncoding.EncodeToString(id.Bytes())
}
