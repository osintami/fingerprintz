package etl

import (
	"strings"
)

type HttpExtractor struct {
	tools  *Toolbox
	source *Source
}

func NewHttpExtractor(tools *Toolbox, source *Source) IExtract {
	return &HttpExtractor{tools: tools, source: source}
}

func (x *HttpExtractor) Extract(job IETLJob) error {
	// NOTE:  -http:// in the configuration file denotes don't download again, which is very
	//   useful for testing ETL sources that may get upset about too many data pulls
	if x.source.URL != "" && !strings.HasPrefix(x.source.URL, "-") {
		url := x.source.URL
		if x.source.ApiKey != "" {
			apiKey := x.tools.Secrets.Find(x.source.ApiKey)
			url = strings.ReplaceAll(x.source.URL, "{key}", apiKey)
		}
		x.tools.Network.DownloadFile(url, job.Info().inputFile)
	}
	return nil
}
