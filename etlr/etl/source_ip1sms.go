// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"encoding/json"

	"github.com/osintami/fingerprintz/common"
)

type IP1SMS struct {
	writer IWriter
}

type IP1SMSInfo struct {
	Phone          string
	LastReportedAt string
	IsDisposable   bool
}
type IP1SMSRow struct {
	Key    string     `json:"Key"`
	Result IP1SMSInfo `json:"Result"`
}

func NewIP1SMS(writer IWriter) ITransform {
	return &IP1SMS{writer: writer}
}

func (x *IP1SMS) Transform(job IETLJob) error {

	burners := make(map[string]string)
	err := common.LoadJson(job.Info().inputFile, &burners)

	for k, v := range burners {
		row := &IP1SMSRow{}
		row.Key = k
		row.Result.Phone = k
		row.Result.LastReportedAt = v
		row.Result.IsDisposable = true
		raw, _ := json.Marshal(row)
		x.writer.Insert(row.Key, raw)
	}

	job.Tools().Items["ip1sms.IsDisposable"] = Item{
		Item:        "phone/ip1sms/IsDisposable",
		Enabled:     true,
		GJSON:       "Result.IsDisposable",
		Description: "Phone is associated with burner/disposable phone services.",
		Type:        common.Boolean.String()}

	job.Tools().Items["ip1sms.lastReportedAt"] = Item{
		Item:        "phone/ip1sms/lastReportedAt",
		Enabled:     true,
		GJSON:       "Result.LastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
