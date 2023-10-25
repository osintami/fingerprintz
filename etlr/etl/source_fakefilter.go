// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"encoding/json"
	"time"

	"github.com/osintami/fingerprintz/common"
)

type Fakefilter struct {
	writer IWriter
}

type FakeFilterInfo struct {
	Domain         string
	IsFake         bool
	LastReportedAt string
}
type FakeFilterRow struct {
	Key    string         `json:"Key"`
	Result FakeFilterInfo `json:"Result"`
}

func NewFakefilter(writer IWriter) ITransform {
	return &Fakefilter{writer: writer}
}

func (x *Fakefilter) Transform(job IETLJob) error {

	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		row := &FakeFilterRow{}
		row.Key = values[0]
		row.Result.Domain = values[0]
		row.Result.IsFake = true
		row.Result.LastReportedAt = lastReportedAt
		raw, _ := json.Marshal(row)
		x.writer.Insert(row.Key, raw)
		return nil
	})

	job.Tools().Items["domain.fakefilter.IsNefariusDomain"] = Item{
		Item:        "domain/fakefilter/IsFake",
		Enabled:     true,
		GJSON:       "Result.IsFake",
		Description: "Domain is associated with fake/disposable emails.",
		Type:        common.Boolean.String()}
	job.Tools().Items["domain.fakefilter.lastReportedAt"] = Item{
		Item:        "domain/fakefilter/lastReportedAt",
		Enabled:     true,
		GJSON:       "Result.LastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
