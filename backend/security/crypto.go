package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/mr-tron/base58"
)

type KeyPair struct {
	PrivateKey ed25519.PrivateKey `json:"private_key"`
	PublicKey  ed25519.PublicKey  `json:"public_key"`
}

type CryptoManager struct {
	keyPair   *KeyPair
	keyPath   string
	nickname  string
}

type IdentityData struct {
	Nickname  string `json:"nickname"`
	PublicKey string `json:"public_key"`
}

func NewCryptoManager(keyPath string) *CryptoManager {
	return &CryptoManager{
		keyPath: keyPath,
	}
}

func (cm *CryptoManager) GenerateKeyPair() error {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	
	cm.keyPair = &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
	return nil
}

func (cm *CryptoManager) LoadOrGenerateKeys(nickname string) error {
	cm.nickname = nickname
	
	if err := cm.loadKeys(); err != nil {
		if err := cm.GenerateKeyPair(); err != nil {
			return err
		}
		return cm.saveKeys()
	}
	return nil
}

func (cm *CryptoManager) loadKeys() error {
	if _, err := os.Stat(cm.keyPath); os.IsNotExist(err) {
		return errors.New("key file does not exist")
	}
	
	data, err := ioutil.ReadFile(cm.keyPath)
	if err != nil {
		return err
	}
	
	var identity IdentityData
	if err := json.Unmarshal(data, &identity); err != nil {
		return err
	}
	
	cm.nickname = identity.Nickname
	publicKeyBytes, err := hex.DecodeString(identity.PublicKey)
	if err != nil {
		return err
	}
	
	privateKeyPath := cm.keyPath + ".private"
	privateKeyData, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}
	
	privateKeyBytes, err := hex.DecodeString(string(privateKeyData))
	if err != nil {
		return err
	}
	
	cm.keyPair = &KeyPair{
		PrivateKey: ed25519.PrivateKey(privateKeyBytes),
		PublicKey:  ed25519.PublicKey(publicKeyBytes),
	}
	
	return nil
}

func (cm *CryptoManager) saveKeys() error {
	if err := os.MkdirAll(filepath.Dir(cm.keyPath), 0700); err != nil {
		return err
	}
	
	identity := IdentityData{
		Nickname:  cm.nickname,
		PublicKey: hex.EncodeToString(cm.keyPair.PublicKey),
	}
	
	data, err := json.MarshalIndent(identity, "", "  ")
	if err != nil {
		return err
	}
	
	if err := ioutil.WriteFile(cm.keyPath, data, 0644); err != nil {
		return err
	}
	
	privateKeyPath := cm.keyPath + ".private"
	privateKeyHex := hex.EncodeToString(cm.keyPair.PrivateKey)
	return ioutil.WriteFile(privateKeyPath, []byte(privateKeyHex), 0600)
}

func (cm *CryptoManager) SignMessage(message []byte) []byte {
	if cm.keyPair == nil {
		return nil
	}
	return ed25519.Sign(cm.keyPair.PrivateKey, message)
}

func (cm *CryptoManager) VerifySignature(message, signature []byte, publicKey ed25519.PublicKey) bool {
	return ed25519.Verify(publicKey, message, signature)
}

func (cm *CryptoManager) GetPublicKey() ed25519.PublicKey {
	if cm.keyPair == nil {
		return nil
	}
	return cm.keyPair.PublicKey
}

func (cm *CryptoManager) GetPrivateKey() ed25519.PrivateKey {
	if cm.keyPair == nil {
		return nil
	}
	return cm.keyPair.PrivateKey
}

func (cm *CryptoManager) GetPublicKeyFingerprint() string {
	if cm.keyPair == nil {
		return ""
	}
	hash := sha256.Sum256(cm.keyPair.PublicKey)
	return hex.EncodeToString(hash[:8])
}

func (cm *CryptoManager) GetPublicKeyBase58() string {
	if cm.keyPair == nil {
		return ""
	}
	return base58.Encode(cm.keyPair.PublicKey)
}

func (cm *CryptoManager) GetNickname() string {
	return cm.nickname
}

func (cm *CryptoManager) EncryptAES(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (cm *CryptoManager) DecryptAES(encryptedData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("encrypted data too short")
	}
	
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
} 