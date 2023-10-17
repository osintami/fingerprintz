// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/log"
)

type MaxmindWriter struct {
	tree *mmdbwriter.Tree
}

func NewMaxmindWriter(name string) *MaxmindWriter {
	desc := make(map[string]string)
	desc["en"] = "OSINTAMI"
	opts := mmdbwriter.Options{
		IPVersion:    6,
		RecordSize:   28,
		DatabaseType: name,
		Languages:    []string{"en"},
		Description:  desc}

	tree, err := mmdbwriter.New(opts)
	if err != nil || name == "" {
		log.Error().Err(err).Str("component", "mmdb").Str("name", name).Msg("create tree")
		return nil
	}
	return &MaxmindWriter{tree: tree}
}

func (x *MaxmindWriter) Insert(cidr *net.IPNet, entry mmdbtype.Map) error {
	return x.tree.InsertFunc(cidr, inserter.TopLevelMergeWith(entry))
}

func (x *MaxmindWriter) Close(mmdbFile string) error {
	fh, err := os.Create(mmdbFile)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = x.tree.WriteTo(fh)
	return err
}
