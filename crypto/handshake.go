package crypto

import (
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"io"
)

const publicKeySize = 65

func ClientHandshake(rw io.ReadWriter) (*Session, error) {
	curve := ecdh.P256()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	clientPub := privateKey.PublicKey().Bytes()
	if err := writeExact(rw, clientPub); err != nil {
		return nil, err
	}

	serverPubBytes, err := readExact(rw, publicKeySize)
	if err != nil {
		return nil, err
	}

	serverPub, err := curve.NewPublicKey(serverPubBytes)
	if err != nil {
		return nil, err
	}

	shared, err := privateKey.ECDH(serverPub)
	if err != nil {
		return nil, err
	}

	key, err := deriveKey(shared, clientPub, serverPubBytes)
	if err != nil {
		return nil, err
	}

	return NewSession(key)
}

func ServerHandshake(rw io.ReadWriter) (*Session, error) {
	curve := ecdh.P256()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	clientPubBytes, err := readExact(rw, publicKeySize)
	if err != nil {
		return nil, err
	}

	serverPub := privateKey.PublicKey().Bytes()
	if err := writeExact(rw, serverPub); err != nil {
		return nil, err
	}

	clientPub, err := curve.NewPublicKey(clientPubBytes)
	if err != nil {
		return nil, err
	}

	shared, err := privateKey.ECDH(clientPub)
	if err != nil {
		return nil, err
	}

	key, err := deriveKey(shared, clientPubBytes, serverPub)
	if err != nil {
		return nil, err
	}

	return NewSession(key)
}

func deriveKey(shared, clientPub, serverPub []byte) ([]byte, error) {
	salt := append(append([]byte{}, clientPub...), serverPub...)
	return hkdf.Key(sha256.New, shared, salt, "sbtp-handshake", 32)
}

func writeExact(w io.Writer, data []byte) error {
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(data)))
	if _, err := w.Write(length); err != nil {
		return err
	}
	_, err := w.Write(data)
	return err
}

func readExact(r io.Reader, expected int) ([]byte, error) {
	length := make([]byte, 2)
	if _, err := io.ReadFull(r, length); err != nil {
		return nil, err
	}
	n := int(binary.BigEndian.Uint16(length))
	if n != expected {
		return nil, io.ErrUnexpectedEOF
	}
	data := make([]byte, n)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}
