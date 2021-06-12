//Package govapid is a micro-package to generate VAPID keys which are required for web-push
//Voluntary Application Server Identification (VAPID) for Web Push https://datatracker.ietf.org/doc/rfc8292/
package govapid

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
)

//VAPIDKeys contains the public and private VAPID keys
type VAPIDKeys struct {
	Public  string
	Private string
}

//GenerateVAPID will generate public and private VAPID keys
func GenerateVAPID() (VAPIDKeys, error) {

	curve := elliptic.P256()
	privateKey, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return VAPIDKeys{}, err
	}

	publicKey := elliptic.Marshal(curve, x, y)

	privKey := base64.RawURLEncoding.EncodeToString(privateKey)
	pubKey := base64.RawURLEncoding.EncodeToString(publicKey)

	return VAPIDKeys{pubKey, privKey}, nil
}
