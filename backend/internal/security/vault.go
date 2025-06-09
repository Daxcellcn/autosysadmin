// backend/internal/security/vault.go
package security

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"
)

type SecretVault interface {
	StoreSecret(ctx context.Context, agentID, key string, value []byte) error
	RetrieveSecret(ctx context.Context, agentID, key string) ([]byte, error)
	DeleteSecret(ctx context.Context, agentID, key string) error
	ListSecrets(ctx context.Context, agentID string) ([]string, error)
	RotateEncryptionKey() error
}

type secretVault struct {
	secrets      map[string]map[string][]byte // agentID -> key -> encrypted value
	encryptionKey []byte
	mu           sync.RWMutex
}

func NewVault(encryptionKey string) SecretVault {
	return &secretVault{
		secrets:      make(map[string]map[string][]byte),
		encryptionKey: []byte(encryptionKey),
	}
}

func (v *secretVault) StoreSecret(ctx context.Context, agentID, key string, value []byte) error {
	encrypted, err := v.encrypt(value)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	if _, exists := v.secrets[agentID]; !exists {
		v.secrets[agentID] = make(map[string][]byte)
	}

	v.secrets[agentID][key] = encrypted
	return nil
}

func (v *secretVault) RetrieveSecret(ctx context.Context, agentID, key string) ([]byte, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	agentSecrets, exists := v.secrets[agentID]
	if !exists {
		return nil, fmt.Errorf("no secrets found for agent %s", agentID)
	}

	encrypted, exists := agentSecrets[key]
	if !exists {
		return nil, fmt.Errorf("secret %s not found for agent %s", key, agentID)
	}

	return v.decrypt(encrypted)
}

func (v *secretVault) DeleteSecret(ctx context.Context, agentID, key string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	agentSecrets, exists := v.secrets[agentID]
	if !exists {
		return fmt.Errorf("no secrets found for agent %s", agentID)
	}

	if _, exists := agentSecrets[key]; !exists {
		return fmt.Errorf("secret %s not found for agent %s", key, agentID)
	}

	delete(agentSecrets, key)
	return nil
}

func (v *secretVault) ListSecrets(ctx context.Context, agentID string) ([]string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	agentSecrets, exists := v.secrets[agentID]
	if !exists {
		return nil, fmt.Errorf("no secrets found for agent %s", agentID)
	}

	keys := make([]string, 0, len(agentSecrets))
	for key := range agentSecrets {
		keys = append(keys, key)
	}

	return keys, nil
}

func (v *secretVault) RotateEncryptionKey() error {
	newKey := make([]byte, 32)
	if _, err := rand.Read(newKey); err != nil {
		return fmt.Errorf("failed to generate new encryption key: %w", err)
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	// Re-encrypt all secrets with the new key
	for agentID, agentSecrets := range v.secrets {
		for key, encrypted := range agentSecrets {
			decrypted, err := v.decrypt(encrypted)
			if err != nil {
				return fmt.Errorf("failed to decrypt secret during key rotation: %w", err)
			}

			v.encryptionKey = newKey
			newEncrypted, err := v.encrypt(decrypted)
			if err != nil {
				return fmt.Errorf("failed to re-encrypt secret during key rotation: %w", err)
			}

			v.secrets[agentID][key] = newEncrypted
		}
	}

	return nil
}

func (v *secretVault) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(v.encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (v *secretVault) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(v.encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}