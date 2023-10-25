// Copyright © 2023 OSINTAMI. This is not yours.
package etl

type FileDBWriter struct {
	outputType string
}

func NewFileDBWriter(outputType string) IWriter {
	return &FileDBWriter{outputType: outputType}
}

func (x *FileDBWriter) Type() string {
	return x.outputType
}

func (x *FileDBWriter) Create(fileName string) error {
	return nil
}

func (x *FileDBWriter) Load(job IETLJob) error {
	return job.Tools().FileSystem.Copy(job.Info().inputFile, job.Info().snapshotFile)
}

func (x *FileDBWriter) Insert(key interface{}, value interface{}) error {
	return nil
}
