package outer

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"

	"qoobing.com/gomod/log"
)

// KeysClass struct
type KeysClass struct {
	Use string `json:"use"`
	Alg string `json:"alg"`
	Id  string `json:"id"`
	Key string `json:"key"`
}

// request response key
type reqResponse struct {
	ErrCode int         `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Value   []KeysClass `json:"value"`
}

// request arguments
type requestArgs struct {
	Name string `json:"name"`
}

type ecdsaKey ecdsa.PublicKey

// NewEcdsaKey is to generate ecdsa key from public key interface
func NewEcdsaKey(rootURL string, scene string) ([]string, []ecdsaKey, error) {
	var keysId []string
	var ecdsaKeys []ecdsaKey
	keyAddr := rootURL + "/openapi/get_public_keys"
	keys, err := fetchPublicKeyFromApi(keyAddr, scene)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch public key from api error %v", err)
	}
	for _, key := range keys.Value {
		keysId = append(keysId, key.Id)
		publicKey, err := stringToEcdsaPublicKey(key.Key)
		if err != nil {
			return nil, nil, fmt.Errorf("string to ecdsa public key error %v", err)
		}
		ecdsaKeys = append(ecdsaKeys, ecdsaKey(*publicKey))
	}
	return keysId, ecdsaKeys, nil
}

// Sign sign a data, returns encoded signature or error
func (ek ecdsaKey) Sign(data []byte) ([]byte, error) {
	log.Panicf("cannot sign due to ecdsaKey only have public key")
	return nil, nil
}

// Verify verify sign, returns nil if signature is valid
func (ek ecdsaKey) Verify(data []byte, sig []byte) error {
	hash := sha256.Sum256(data)
	publicKey := (*ecdsa.PublicKey)(&ek)
	if valid := ecdsa.VerifyASN1(publicKey, hash[:], sig); !valid {
		return errors.New("invalid sig")
	}
	return nil
}

// GetPublicKey get the verify publickey
func (ek ecdsaKey) GetPublicKey() interface{} {
	publicKey := (*ecdsa.PublicKey)(&ek)
	return publicKey
}

// Encrypt encrypt data by kms
func (ek ecdsaKey) Encrypt(plaintext []byte) (crypttext []byte, err error) {
	panic("unsuported method")
}

// Decrypt decrypt data by kms
func (ek ecdsaKey) Decrypt(crypttext []byte) (plaintext []byte, err error) {
	panic("unsuported method")
}

func fetchPublicKeyFromApi(ipAddr string, name string) (keys reqResponse, err error) {
	requestarg := requestArgs{Name: name}
	jsonData, err := json.Marshal(requestarg)
	if err != nil {
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", ipAddr, bytes.NewBuffer(jsonData))
	if err != nil {
		return reqResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return reqResponse{}, err
	}
	defer resp.Body.Close()
	var publicKeys reqResponse
	err = json.NewDecoder(resp.Body).Decode(&publicKeys)
	if err != nil {
		return reqResponse{}, err
	}
	return publicKeys, nil
}
func stringToEcdsaPublicKey(pubkeyStr string) (*ecdsa.PublicKey, error) {
	keyStr := fmt.Sprintf("-----BEGIN ECDSA_PUBLIC_KEY-----\n%s\n-----END ECDSA_PUBLIC_KEY-----\n", pubkeyStr)
	pemBlock, _ := pem.Decode([]byte(keyStr))
	if pemBlock == nil {
		return nil, fmt.Errorf("err")
	}
	pubKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to convert public key to ECDSA public key")
	}
	return ecdsaPubKey, nil
}
