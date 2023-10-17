// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"net"
	"os"
	"sync"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/log"
)

type MMDBWriter struct {
	tree  *mmdbwriter.Tree
	mutex sync.Mutex
}

func NewMMDBWriter() IWriter {
	return &MMDBWriter{mutex: sync.Mutex{}}
}

func (x *MMDBWriter) Type() string {
	return "mmdb"
}

func (x *MMDBWriter) Create(mmdbName string) error {
	desc := make(map[string]string)
	desc["en"] = mmdbName
	opts := mmdbwriter.Options{
		IPVersion:    6,
		RecordSize:   28,
		DatabaseType: mmdbName,
		Languages:    []string{"en"},
		Description:  desc}

	// create an mmdb tree writer
	tree, err := mmdbwriter.New(opts)
	if err != nil {
		log.Error().Err(err).Str("component", "mmdb").Str("name", mmdbName).Msg("failed to create writer")
		return err
	}
	x.tree = tree
	return nil
}

func (x *MMDBWriter) Load(job IETLJob) error {
	mmdbFile := job.Info().snapshotFile
	fh, err := os.Create(mmdbFile)
	if err != nil {
		log.Error().Err(err).Str("component", "mmdb").Str("file", mmdbFile).Msg("create output mmdb database")
		return err
	}
	defer fh.Close()
	_, err = x.tree.WriteTo(fh)
	if err != nil {
		log.Error().Err(err).Str("file", mmdbFile).Msg("write output mmdb database")
		return err
	}
	return nil
}

func (x *MMDBWriter) Insert(key interface{}, value interface{}) error {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	// NOTE:  mmdb tree inserts are not thread safe
	return x.tree.InsertFunc(key.(*net.IPNet), inserter.TopLevelMergeWith(value.(mmdbtype.Map)))
}
