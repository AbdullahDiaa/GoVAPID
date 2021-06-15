//Package govapid is a micro-package to generate VAPID public and private keys and VAPID authentication headers, required for sending web push notifications.
//The library only supports VAPID-draft-02+ specification.
// https://datatracker.ietf.org/doc/html/draft-ietf-webpush-vapid-02
package govapid

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	_ "crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
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

// validateVAPIDKeys will validate the length and encoding of VAPID keys
func validateVAPIDKeys(keys VAPIDKeys) error {
	if len(keys.Public) != 87 {
		return errors.New("Invalid Public key length")
	}

	if len(keys.Private) != 43 {
		return errors.New("Invalid Private key length")
	}

	_, err := base64.RawURLEncoding.DecodeString(keys.Private)
	if err != nil {
		return errors.New("Invalid Private key")
	}

	_, err = base64.RawURLEncoding.DecodeString(keys.Public)
	if err != nil {
		return errors.New("Invalid Public key")
	}
	return nil
}

//verifyClaims will verify the claims of JWT string
func verifyClaims(claims map[string]interface{}) error {
	//Validate claims
	// sub: The “Subscriber” a mailto link for the administrative contact for this feed.
	// It’s best if this email is not a personal email address,
	// but rather a group email so that if a person leaves an organization,
	// is unavailable for an extended period, or otherwise can’t respond, someone else on the list can.
	if _, ok := claims["sub"]; ok {
		if !(strings.HasPrefix(claims["sub"].(string), "mailto:")) && !(strings.HasPrefix(claims["sub"].(string), "https://")) {
			return errors.New("“Subscriber” claim (sub) is invalid, it should be an email or contact URL")
		}
	}

	//exp : “Expires” this is an integer that is the date and time that this VAPID header should remain valid until.
	// It doesn’t reflect how long your VAPID signature key should be valid, just this specific update.
	// It can be no more than 24 hours
	if _, ok := claims["exp"]; ok {
		now := time.Now().Unix()
		tomorrow := time.Now().Add(24 * time.Hour).Unix()
		if now > claims["exp"].(int64) {
			return errors.New("Expiry claim (exp) already expired")
		}
		if claims["exp"].(int64) > tomorrow {
			return errors.New("Expiry claim (exp) maximum value is 24 hours")
		}
	}
	return nil
}

func generateJWTSignature(keys VAPIDKeys, JWTInfoAndData string) (string, error) {
	// Signature is the third part of the token, which includes the data above signed with the private key
	// Preparing ecdsa.PrivateKey for signing
	privKeyDecoded, err := base64.RawURLEncoding.DecodeString(keys.Private)
	if err != nil {
		return "", errors.New("Invalid VAPID private key string, cannot decode it")
	}

	curve := elliptic.P256()
	px, py := curve.ScalarMult(
		curve.Params().Gx,
		curve.Params().Gy,
		privKeyDecoded,
	)

	pubKey := ecdsa.PublicKey{
		Curve: curve,
		X:     px,
		Y:     py,
	}

	// Private key
	d := &big.Int{}
	d.SetBytes(privKeyDecoded)

	privKey := &ecdsa.PrivateKey{
		PublicKey: pubKey,
		D:         d,
	}

	// Get the key
	hash := crypto.SHA256
	hasher := hash.New()
	hasher.Write([]byte(JWTInfoAndData))

	// Sign JWTInfo and JWTData using the private key
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hasher.Sum(nil))
	if err != nil {
		return "", errors.New("Err singing data")
	}

	curveBits := privKey.Curve.Params().BitSize

	if curveBits != 256 {
		return "", errors.New("curveBits should be 256")
	}

	keyBytes := curveBits / 8
	if curveBits%8 > 0 {
		keyBytes++
	}

	rBytes := r.Bytes()
	rBytesPadded := make([]byte, keyBytes)
	copy(rBytesPadded[keyBytes-len(rBytes):], rBytes)

	sBytes := s.Bytes()
	sBytesPadded := make([]byte, keyBytes)
	copy(sBytesPadded[keyBytes-len(sBytes):], sBytes)

	out := append(rBytesPadded, sBytesPadded...)

	return "." + strings.TrimRight(base64.URLEncoding.EncodeToString(out), "="), nil
}

//GenerateVAPIDAuth will generate Authorization header for web push notifications
func GenerateVAPIDAuth(keys VAPIDKeys, claims map[string]interface{}) (string, error) {

	//Validate VAPID Keys
	if err := validateVAPIDKeys(keys); err != nil {
		return "", err
	}

	//Verify Claims
	if err := verifyClaims(claims); err != nil {
		return "", err
	}

	// JWTInfo is base64 Encoded {"typ":"JWT","alg":"ES256"} which is the first part of the JWT Token
	JWTInfo := "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9."

	// JWTData is the second part of the token which contains all the claims encoded in base64
	jsonValue, err := json.Marshal(claims)
	if err != nil {
		return "", errors.New("Marshaling Claims JSON failed" + err.Error())
	}
	JWTData := strings.TrimRight(base64.URLEncoding.EncodeToString(jsonValue), "=")

	JWTSignature, err := generateJWTSignature(keys, JWTInfo+JWTData)
	if err != nil {
		return "", err
	}

	//Compose the JWT Token string
	JWTString := JWTInfo + JWTData + JWTSignature

	// Construct the VAPID header
	VAPIDAuth := fmt.Sprintf(
		"vapid t=%s, k=%s",
		JWTString,
		keys.Public,
	)

	return VAPIDAuth, nil
}
