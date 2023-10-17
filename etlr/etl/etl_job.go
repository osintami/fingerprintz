// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/osintami/fingerprintz/log"
	"golang.org/x/exp/maps"
)

type ETLJobInfo struct {
	inputFile    string
	outputFile   string
	schemaFile   string
	workingPath  string
	snapshotFile string
	tmpPath      string
	snapshotName string
}

type IETLJob interface {
	Refresh() error
	Extract() error
	Transform() error
	Load() error
	Publish() error
	Info() *ETLJobInfo
	Source() *Source
	Tools() *Toolbox
}

type ETLJob struct {
	tools     *Toolbox
	source    *Source
	info      *ETLJobInfo
	writer    IWriter
	extract   IExtract
	transform ITransform
	load      ILoad
}

type IExtract interface {
	Extract(IETLJob) error
}

type ILoad interface {
	Load(IETLJob) error
}

type ITransform interface {
	Transform(IETLJob) error
}

func NewETLJob(tools *Toolbox, source *Source, dataPath string, writer IWriter, extract IExtract, transform ITransform, load ILoad) *ETLJob {

	tmpPath := "/tmp/"

	tm := time.Now()
	snapshotName := fmt.Sprintf("%s_%d_%02d_%02d.%s", source.Name, tm.Year(), tm.Month(), tm.Day(), source.OutputType)

	info := &ETLJobInfo{
		inputFile: fmt.Sprintf("%s/%s.%s",
			tmpPath+source.Name,
			source.Name,
			source.InputType),
		outputFile:   dataPath + source.Name + "." + source.OutputType,
		schemaFile:   dataPath + source.Name + ".json",
		workingPath:  tmpPath + source.Name + "/",
		snapshotFile: tmpPath + source.Name + "/" + snapshotName,
		tmpPath:      tmpPath,
		snapshotName: snapshotName,
	}

	return &ETLJob{
		tools:     tools,
		source:    source,
		info:      info,
		writer:    writer,
		extract:   extract,
		transform: transform,
		load:      load}

}

func (x *ETLJob) Source() *Source {
	return x.source
}

func (x *ETLJob) Tools() *Toolbox {
	return x.tools
}

func (x *ETLJob) Info() *ETLJobInfo {
	return x.info
}

func (x *ETLJob) Refresh() error {

	log.Info().Str("component", "etlr").Str("state", "start").Str("vendor", x.source.Name).Msg("refresh")

	if err := x.cleanupETL(); err != nil {
		return err
	}
	if err := x.prepareETL(); err != nil {
		return err
	}
	if err := x.Extract(); err != nil {
		return err
	}
	if err := x.Transform(); err != nil {
		return err
	}
	if err := x.Load(); err != nil {
		return err
	}
	err := x.Publish()

	log.Info().Str("component", "etlr").Str("state", "finish").Str("vendor", x.source.Name).Msg("refresh")
	return err
}

func (x *ETLJob) Extract() error {
	return x.extract.Extract(x)
}

func (x *ETLJob) Transform() error {
	return x.transform.Transform(x)
}

func (x *ETLJob) Load() error {
	return x.load.Load(x)
}

func (x *ETLJob) Publish() error {
	err1 := x.publishDatums()
	err2 := x.publishSchema()
	if err1 != nil {
		return err1
	}
	return err2
}

func (x *ETLJob) prepareETL() error {
	err := os.MkdirAll(x.info.workingPath, 01777)
	if err != nil {
		log.Error().Err(err).Str("component", "etlr").Str("vendor", x.source.Name).Msg("create working directory")
		return err
	}
	switch x.source.OutputType {
	case "mmdb":
		err = x.writer.Create(x.source.Name)
	// NOTE:  csv output works, but very little is instrumented, on hold for now
	// case "csv":
	// 	err = x.writer.Create(x.info.workingPath + x.info.snapshotName)
	case "fast":
		err = x.writer.Create(x.info.workingPath + x.info.snapshotName)
	}
	if err != nil {
		log.Error().Err(err).Str("component", "etlr").Str("vendor", x.source.Name).Msg("create/load mmdb database")
		return err
	}
	return nil
}

func (x *ETLJob) publishDatums() error {
	inFile := x.info.snapshotFile
	outFile := x.info.outputFile

	err := x.tools.FileSystem.Move(inFile, outFile)
	if err != nil {
		log.Error().Err(err).Str("component", "etl").Str("vendor", x.source.Name).Str("inFile", inFile).Str("outFile", outFile).Msg("data publish")
		return err
	}

	uname, err := user.Current()
	if err == nil {
		x.tools.FileSystem.Chown(outFile, uname.Username)
	}
	return nil
}

func (x *ETLJob) publishSchema() error {
	if len(x.tools.Items) == 0 {
		return ErrEmptySchema
	}
	// create data dictionary for this run
	data, _ := json.MarshalIndent(maps.Values(x.tools.Items), "", "    ")

	outFile := x.info.schemaFile
	err := os.WriteFile(outFile, data, 0644)
	if err != nil {
		log.Error().Err(err).Str("component", "etl").Str("vendor", x.source.Name).Str("outFile", outFile).Msg("schema publish")
		return err
	}

	uname, err := user.Current()
	if err == nil {
		x.tools.FileSystem.Chown(outFile, uname.Username)
	}
	return nil
}

func (x *ETLJob) cleanupETL() error {
	for k := range x.tools.Items {
		delete(x.tools.Items, k)
	}
	return os.RemoveAll(x.info.workingPath)
}
