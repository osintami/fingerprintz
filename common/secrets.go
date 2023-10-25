// Copyright © 2023 OSINTAMI. This is not yours.
package common

import "os"

type ISecrets interface {
	Set(key, value string)
	Find(sourceName string) string
}

type Secrets struct {
	secrets map[string]string
}

func NewSecrets(keys []string) *Secrets {
	secrets := make(map[string]string)
	for _, key := range keys {
		secrets[key] = os.Getenv(key)
	}
	return &Secrets{secrets: secrets}
}

func (x *Secrets) Set(key, value string) {
	x.secrets[key] = value
}

func (x *Secrets) Find(sourceName string) string {
	return x.secrets[sourceName]
}
