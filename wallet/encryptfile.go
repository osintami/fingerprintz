// Copyright © 2024 OSINTAMI. This is not yours.
package main

import (
	"encoding/json"
	"os"

	"github.com/osintami/fingerprintz/log"
)

type FileEncrypt struct {
	encrypt ICrypto
}

func NewFileEncrypt(encrypt ICrypto) *FileEncrypt {
	return &FileEncrypt{encrypt: encrypt}
}

func (x *FileEncrypt) ReadFile(fileName string) ([]byte, error) {
	contents, err := os.ReadFile(fileName)
	if err != nil {
		log.Error().Err(err).Str("component", "encrypt").Msg("file open")
		return nil, err
	}
	data, err := x.encrypt.Decrypt(contents)
	if err != nil {
		log.Error().Err(err).Str("component", "encrypt").Msg("decrypt")
		return nil, err
	}
	return data, nil
}

// TODO:  make me generic
func (x *FileEncrypt) ToJSON(fileName string, secrets *[]*Secret) error {
	data, err := x.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, secrets)
	if err != nil {
		log.Error().Err(err).Str("component", "encrypt").Msg("json unmarshal")
		return err
	}
	return nil
}

func (x *FileEncrypt) FromJSON(fileName string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Error().Err(err).Str("component", "encrypt").Msg("json marshal")
		return err
	}
	return x.WriteFile(fileName, data)
}

func (x *FileEncrypt) WriteFile(fileName string, data []byte) error {
	data, err := x.encrypt.Encrypt(data)
	if err != nil {
		log.Error().Err(err).Str("component", "encrypt").Msg("encrypt data")
		return err
	}
	err = os.WriteFile(fileName, data, 0777)
	if err != nil {
		log.Error().Err(err).Str("component", "encrypt").Msg("file write")
		return err
	}
	return nil
}
