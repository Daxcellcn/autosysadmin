// backend/internal/security/ssh.go
package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

type SSHKeyManager interface {
	GenerateKeyPair(bits int) (privateKey, publicKey string, err error)
	RotateKeys(agentID string) error
	GetAuthorizedKeys(agentID string) ([]string, error)
	AddAuthorizedKey(agentID, publicKey string) error
	RemoveAuthorizedKey(agentID, keyID string) error
}

type sshKeyManager struct {
	keys map[string][]SSHKey // agentID -> keys
	mu   sync.RWMutex
}

type SSHKey struct {
	ID         string    `json:"id"`
	PublicKey  string    `json:"public_key"`
	PrivateKey string    `json:"private_key,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

func NewSSHKeyManager() SSHKeyManager {
	return &sshKeyManager{
		keys: make(map[string][]SSHKey),
	}
}

func (m *sshKeyManager) GenerateKeyPair(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Encode private key to PEM
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyPEM)

	// Generate public key
	pubKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate public key: %w", err)
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(pubKey)

	return string(privateKeyBytes), string(publicKeyBytes), nil
}

func (m *sshKeyManager) RotateKeys(agentID string) error {
	// Generate new key pair
	privateKey, publicKey, err := m.GenerateKeyPair(2048)
	if err != nil {
		return err
	}

	// Store the new key
	key := SSHKey{
		ID:         fmt.Sprintf("key-%s-%d", agentID, time.Now().Unix()),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().AddDate(1, 0, 0), // 1 year
	}

	m.mu.Lock()
	m.keys[agentID] = append(m.keys[agentID], key)
	m.mu.Unlock()

	return nil
}

func (m *sshKeyManager) GetAuthorizedKeys(agentID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys, exists := m.keys[agentID]
	if !exists {
		return nil, fmt.Errorf("no keys found for agent %s", agentID)
	}

	var publicKeys []string
	for _, key := range keys {
		publicKeys = append(publicKeys, key.PublicKey)
	}

	return publicKeys, nil
}

func (m *sshKeyManager) AddAuthorizedKey(agentID, publicKey string) error {
	key := SSHKey{
		ID:         fmt.Sprintf("key-%s-%d", agentID, time.Now().Unix()),
		PublicKey:  publicKey,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().AddDate(1, 0, 0), // 1 year
	}

	m.mu.Lock()
	m.keys[agentID] = append(m.keys[agentID], key)
	m.mu.Unlock()

	return nil
}

func (m *sshKeyManager) RemoveAuthorizedKey(agentID, keyID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	keys, exists := m.keys[agentID]
	if !exists {
		return fmt.Errorf("no keys found for agent %s", agentID)
	}

	var updatedKeys []SSHKey
	for _, key := range keys {
		if key.ID != keyID {
			updatedKeys = append(updatedKeys, key)
		}
	}

	m.keys[agentID] = updatedKeys
	return nil
}