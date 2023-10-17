// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"

	"github.com/osintami/fingerprintz/log"
)

type FileExtractor struct {
	fileName string
}

func NewFileExtractor(fileName string) IExtract {
	return &FileExtractor{fileName: fileName}
}

func (x *FileExtractor) Extract(job IETLJob) error {
	data, err := os.ReadFile(x.fileName)
	if err != nil {
		log.Error().Err(err).Str("component", "file extractor").Msg("read file")
		return err
	}
	return os.WriteFile(job.Info().inputFile, data, 0644)
}
