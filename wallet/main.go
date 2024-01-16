// Copyright © 2024 OSINTAMI. This is not yours.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/osintami/fingerprintz/log"
	"golang.org/x/term"
)

type Secret struct {
	Key   string
	Site  string `json:",omitempty"`
	User  string
	Pass  string
	Email string `json:",omitempty"`
	Phone string `json:",omitempty"`
	Code  string `json:",omitempty"`
}

func main() {

	var getKey, setKey, delKey string

	flag.StringVar(&getKey, "get", "", "get account details")
	flag.StringVar(&setKey, "set", "", "set/edit account details")
	flag.StringVar(&delKey, "del", "", "delete account details")
	showKey := flag.Bool("show", false, "show keys")

	var encryptFile, decryptFile, outputFile string
	flag.StringVar(&encryptFile, "encrypt", "", "file to encrypt")
	flag.StringVar(&decryptFile, "decrypt", "", "file to decrypt")
	flag.StringVar(&outputFile, "out", "", "file to output")

	flag.Parse()

	if !*showKey && (encryptFile != "" || decryptFile != "") && outputFile == "" {
		fmt.Println("To encrypt/decrypt a file you must provide an output file.")
		fmt.Println("./wallet -decrypt my.dat -out my.json")
		return
	}

	encryptKey, err := getPasscode("passcode")
	if err != nil {
		fmt.Println(err)
		return
	}

	crypto := NewEncrypt([]byte(encryptKey))
	fe := NewFileEncrypt(crypto)

	// read/display an account
	if getKey != "" {
		var secrets []*Secret
		err = fe.ToJSON("my.dat", &secrets)
		if err != nil {
			log.Error().Err(err).Str("component", "encrypt").Msg("json unmarshal")
			return
		}

		for _, secret := range secrets {
			if secret.Key == getKey {
				out, _ := json.MarshalIndent(secret, "", "  ")
				fmt.Println(string(out))
				return
			}
		}
		fmt.Println("[ERROR] key not found")
		return
	}

	if *showKey {
		var secrets []*Secret
		err = fe.ToJSON("my.dat", &secrets)
		if err != nil {
			log.Error().Err(err).Str("component", "encrypt").Msg("json unmarshal")
			return
		}
		for _, secret := range secrets {
			fmt.Println(secret.Key)
		}
		return
	}

	if setKey != "" {

		var secrets []*Secret
		err = fe.ToJSON("my.dat", &secrets)
		if err != nil {
			log.Error().Err(err).Str("component", "encrypt").Msg("json unmarshal")
			return
		}

		secret := find(secrets, setKey)
		secret.Key = getInput("key", secret.Key)
		secret.Site = getInput("name", secret.Site)
		secret.User = getInput("user", secret.User)

		fmt.Printf("Password: (%s) ", secret.Pass)
		password, _ := term.ReadPassword(syscall.Stdin)
		fmt.Println()
		pass := string(password)
		if pass == "" {
			pass = secret.Pass
		}
		secret.Pass = pass

		secret.Email = getInput("email", secret.Email)
		secret.Phone = getInput("phone", secret.Phone)

		next, _ := json.MarshalIndent(secret, "", "  ")
		fmt.Println(string(next))

		fmt.Print("Save (Y/N): ")
		prompt := readline()
		fmt.Println()
		if prompt == "Y" || prompt == "y" || prompt == "yes" || prompt == "YES" {
			secrets = append(secrets, secret)
			fe.FromJSON("my.dat", secrets)
		}

		return
	}

	if delKey != "" {

		var secrets []*Secret
		err = fe.ToJSON("my.dat", &secrets)
		if err != nil {
			log.Error().Err(err).Str("component", "encrypt").Msg("json unmarshal")
			return
		}

		secret := find(secrets, delKey)
		out, _ := json.MarshalIndent(secret, "", "  ")
		fmt.Println(string(out))

		fmt.Print("Delete (Y/N): ")
		prompt := readline()
		fmt.Println()
		if prompt == "Y" || prompt == "y" || prompt == "yes" || prompt == "YES" {
			var out []*Secret
			for _, s := range secrets {
				if s.Key != delKey {
					out = append(out, s)
				}
			}
			fe.FromJSON("my.dat", out)
		}

		return
	}

	// encrypt or decrypt accordingly
	if encryptFile != "" && outputFile != "" {

		// double check passphrase/key
		passphrase, err := getPasscode("passcode (retype)")
		if err != nil {
			log.Error().Err(err).Str("component", "encrypt").Msg("passphrase")
			return
		}

		if encryptKey != passphrase {
			log.Error().Err(err).Str("component", "encrypt").Msg("passphrase mismatch")
			return
		}

		contents, err := os.ReadFile(encryptFile)
		if err != nil {
			log.Error().Err(err).Str("component", "encrypt").Msg("encrypt file")
			return
		}
		err = fe.WriteFile(outputFile, contents)
		if err != nil {
			log.Error().Err(err).Str("component", "encrypt").Msg("write file")
			return
		}
		return
	}

	if decryptFile != "" && outputFile != "" {
		data, err := fe.ReadFile(decryptFile)
		if err != nil {
			log.Error().Err(err).Str("component", "encrypt").Msg("decrypt file")
			return
		}
		err = os.WriteFile(outputFile, data, 0660)
		if err != nil {
			log.Error().Err(err).Str("component", "encrypt").Msg("write file")
			return
		}
		return
	}

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func readline() string {
	bio := bufio.NewReader(os.Stdin)
	line, _, err := bio.ReadLine()
	if err != nil {
		fmt.Println(err)
	}
	return string(line)
}

func getPasscode(prompt string) (string, error) {
	// prompt for passphrase with hidden input
	fmt.Printf("%s: ", strings.ToTitle(prompt))
	input, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	key := string(input)
	if len(key) < 32 {
		fmt.Print("NOTE:  you are using a pretty weak pass key")
	}
	fmt.Println()
	return fmt.Sprintf("%032s", key), nil
}

func find(secrets []*Secret, key string) *Secret {
	for _, secret := range secrets {
		if secret.Key == key {
			return secret
		}
	}
	// return an empty secret
	return &Secret{Key: key}
}

func getInput(param, value string) string {
	fmt.Printf("%s (%s): ", strings.ToTitle(param), value)
	key := readline()
	if key == "" {
		return value
	}
	return key
}
