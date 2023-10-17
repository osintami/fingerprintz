// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"gitee.com/golang-module/dongle"
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

func NewAES(key []byte) *Encrypt {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error().Err(err).Str("component", "encrypt").Msg("NewCipher")
		return nil
	}
	return &Encrypt{key: key, block: block}
}

func (x *Encrypt) Encrypt(data []byte) ([]byte, error) {
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	io.ReadFull(rand.Reader, iv)
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

func SHA256ToBase16Lowercase(key string) string {
	sha256 := sha256.New()
	sha256.Write([]byte(key))
	return dongle.Encode.FromBytes(sha256.Sum(nil)).ByBase16().ToString()
}
