// Copyright © 2024 OSINTAMI. This is not yours.
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"github.com/osintami/fingerprintz/log"
)

type ICrypto interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type Encrypt struct {
	key   []byte
	block cipher.Block
}

var ErrInvalidBlockSize = errors.New("invalid block size")

func NewEncrypt(key []byte) *Encrypt {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error().Err(err).Str("component", "encrypt").Msg("NewCipher")
	}
	return &Encrypt{key: key, block: block}
}

func (x *Encrypt) Encrypt(data []byte) ([]byte, error) {
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(x.block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

func (x *Encrypt) Decrypt(data []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return []byte("[]"), ErrInvalidBlockSize
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(x.block, iv)
	stream.XORKeyStream(data, data)
	return data, nil
}
