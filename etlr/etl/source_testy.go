// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type TestSource struct {
	writer IWriter
}

func NewTestSource(writer IWriter) ITransform {
	return &TestSource{writer: writer}
}

func (x *TestSource) Transform(job IETLJob) error {

	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"mock": mmdbtype.Map{
					"test": mmdbtype.Map{
						"isTest":         mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}
			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["mock.test.isTest"] = Item{
		Item:        "ip/mock.test/isTest",
		Enabled:     true,
		GJSON:       "mock.test.isTest",
		Description: "IP belongs to Test.",
		Type:        common.Boolean.String()}
	job.Tools().Items["mock.test.lastReportedAt"] = Item{
		Item:        "ip/mock.test/lastReportedAt",
		Enabled:     true,
		GJSON:       "mock.test.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
