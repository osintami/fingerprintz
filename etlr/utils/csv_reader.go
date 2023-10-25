// Copyright © 2023 OSINTAMI. This is not yours.
package utils

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/osintami/fingerprintz/log"
)

type ICSV interface {
	ProcessFile(rawFile string, sep rune, parseLine func([]string) error) error
}

type CSVReader struct {
}

func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

func (x *CSVReader) ProcessFile(rawFile string, sep rune, parseLine func([]string) error) error {
	fh, err := os.Open(rawFile)
	if err != nil {
		log.Error().Err(err).Str("component", "csv").Str("file", rawFile).Msg("open")
		return err
	}
	defer fh.Close()

	csvReader := csv.NewReader(fh)
	csvReader.TrimLeadingSpace = true
	csvReader.Comment = '#'
	csvReader.Comma = sep

	for {
		recs, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Error().Err(err).Str("component", "csv helper").Str("file", rawFile).Msg("csv reader")
			return err
		}
		for i := range recs {
			recs[i] = strings.Trim(recs[i], " ")
			recs[i] = strings.Trim(recs[i], "\"")
		}
		if recs[0] != "" {
			err := parseLine(recs)
			if err != nil {
				log.Error().Err(err).Str("component", "csv helper").Str("file", rawFile).Str("cidr", strings.Join(recs, ",")).Msg("parse line")
				continue
			}
		}
	}
	return nil
}
