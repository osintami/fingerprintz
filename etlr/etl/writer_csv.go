// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"encoding/csv"
	"os"

	"github.com/osintami/fingerprintz/log"
)

type CSVWriter struct {
	writer *csv.Writer
}

func NewCSVWriter() IWriter {
	return &CSVWriter{}
}

func (x *CSVWriter) Type() string {
	return "csv"
}

func (x *CSVWriter) Create(csvFile string) error {
	fh, err := os.Create(csvFile)
	if err != nil {
		log.Error().Err(err).Str("component", "csv").Str("file", csvFile).Msg("create csv database")
		return err
	}
	x.writer = csv.NewWriter(fh)
	return nil
}

func (x *CSVWriter) Load(job IETLJob) error {
	x.writer.Flush()
	return nil
}

func (x *CSVWriter) Insert(key interface{}, value interface{}) error {
	return x.writer.Write(value.([]string))
}
