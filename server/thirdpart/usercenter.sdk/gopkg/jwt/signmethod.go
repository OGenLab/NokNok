package jwt

import (
	"encoding/base64"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"qoobing.com/gomod/log"
)

// KmsSignMethod kms sign method for jwt
type KmsSignMethod struct {
	Name  string // alg name
	Scene string // kms scene
}

// NewKmsSignMethod new a KmsSignMethod
func NewKmsSignMethod(name, scene string) *KmsSignMethod {
	newKmsSignMethod := KmsSignMethod{
		Name:  name,
		Scene: scene,
	}
	jwt.RegisterSigningMethod(newKmsSignMethod.Alg(), func() jwt.SigningMethod {
		return &newKmsSignMethod
	})
	return &newKmsSignMethod
}

// Alg returns the alg identifier for this method (example: 'HS256')
func (ksm KmsSignMethod) Alg() string {
	return ksm.Name
}

// Verify verify signature, returns nil if signature is valid
func (ksm KmsSignMethod) Verify(signingString, signature string, key interface{}) error {
	log.Debugf("scene=%s, key =%v", ksm.Scene, key)
	var (
		keyid = key.(string)
		data  = []byte(signingString)
	)
	sign, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("invalid signatrue, base64 err: '%s'", err)
	}

	if err = jwtkms.Verify(ksm.Scene, keyid, data, sign); err != nil {
		log.Debugf("Verify failed")
		return fmt.Errorf("invalid signatrue, kms verify err: '%s'", err)
	}
	log.Debugf("Verify success")
	return nil
}

// Sign do signature, Returns encoded signature or error
func (ksm KmsSignMethod) Sign(signingString string, key interface{}) (string, error) {
	log.Debugf("scene=%s, key =%v", ksm.Scene, key)
	var (
		keyid     = key.(string)
		plaintext = []byte(signingString)
	)
	sign, err := jwtkms.Sign(ksm.Scene, keyid, plaintext)
	if err != nil {
		return "", err
	}

	signature := base64.RawURLEncoding.EncodeToString(sign)
	return signature, nil
}
