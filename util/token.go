package util

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string, duration time.Duration) (*Payload, error)
}

type PasetoMaker struct {
	Paseto       *paseto.V2
	SymmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key size: must be exactly %d bytes, got %d bytes", chacha20poly1305.KeySize, len(symmetricKey))
	}

	maker := &PasetoMaker{
		Paseto:       paseto.NewV2(),
		SymmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token, err := maker.Paseto.Encrypt(maker.SymmetricKey, payload, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (maker *PasetoMaker) VerifyToken(token string, duration time.Duration) (*Payload, error) {
	payload := &Payload{}

	err := maker.Paseto.Decrypt(token, maker.SymmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	if time.Now().After(payload.ExpiredAt) {
		return nil, fmt.Errorf("token has expired")
	}

	if time.Now().Before(payload.IssuedAt) {
		return nil, fmt.Errorf("token is not valid yet")
	}

	return payload, nil
}
