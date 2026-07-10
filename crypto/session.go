package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
)

var ErrSessionNotEstablished = errors.New("crypto: session not established")

type Session struct {
	aead     cipher.AEAD
	writeSeq uint64
	readSeq  uint64
}

func NewSession(key []byte) (*Session, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Session{aead: aead}, nil
}

func (s *Session) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, s.aead.NonceSize())
	binary.BigEndian.PutUint64(nonce[4:], s.writeSeq)
	s.writeSeq++

	ciphertext := s.aead.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

func (s *Session) Decrypt(ciphertext []byte) ([]byte, error) {
	nonce := make([]byte, s.aead.NonceSize())
	binary.BigEndian.PutUint64(nonce[4:], s.readSeq)
	s.readSeq++

	return s.aead.Open(nil, nonce, ciphertext, nil)
}
