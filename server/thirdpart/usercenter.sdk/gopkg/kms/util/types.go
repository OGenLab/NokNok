package util

// Key define crypto Sign/Verify/GetPublicKey/Encrypt/Decrypt a key need
type Key interface {
	// Sign sign a data, returns encoded signature or error
	//   @scene scene for public key, one of {patk, uatk, vatk}
	//   @keyid keyid for get privatekey to sign
	//   @data  target data
	Sign(data []byte) ([]byte, error)
	// Verify verify sign, returns nil if signature is valid
	//   @scene scene for public key, one of {patk, uatk, vatk}
	//   @keyid keyid for get publickey to verify token
	//   @data  target data
	Verify(data []byte, sign []byte) error
	// GetPublicKey get the verify publickey
	//   @scene scene for public key, one of {patk, uatk, vatk}
	//   @keyid keyid for get publickey to verify token
	GetPublicKey() interface{}

	// Encrypt encrypt data by key
	Encrypt(plaintext []byte) (crypttext []byte, err error)
	// Decrypt decrypt data by key
	Decrypt(crypttext []byte) (plaintext []byte, err error)
}

// InnerCryptClass inner crypt class
type InnerCryptClass struct {
	Type    string
	Name    string
	Version int
	Keys    map[string]Key
}
