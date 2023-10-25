// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/oschwald/maxminddb-golang"
	"github.com/osintami/fingerprintz/log"
)

var ErrItemNotFound = errors.New("no data found")
var ErrNoDataPresent = errors.New("no data present")
var ErrDataReader = errors.New("reader error")

type IMaxmindReader interface {
	Resync()
	FileName() string
	Lookup(net.IP) (json.RawMessage, error)
}

type MaxmindReader struct {
	fileName string
	reader   *maxminddb.Reader
}

func NewMaxmindReader(fileName string) (IMaxmindReader, error) {
	reader, err := maxminddb.Open(fileName)
	if err != nil {
		log.Error().Err(err).Str("component", "mmdb reader").Str("file", fileName).Msg("open failed")
	}
	return &MaxmindReader{
		fileName: fileName,
		reader:   reader}, err
}

func (x *MaxmindReader) FileName() string {
	return x.fileName
}

func (x *MaxmindReader) Resync() {
	reader, err := maxminddb.Open(x.fileName)
	if err != nil {
		log.Error().Err(err).Str("component", "mmdb reader").Str("file", x.fileName).Msg("open failed")
		return
	}
	if x.reader != nil {
		x.reader.Close()
	}
	x.reader = reader
}

func (x *MaxmindReader) Lookup(ip net.IP) (json.RawMessage, error) {
	var result interface{}
	if x.reader != nil {
		err := x.reader.Lookup(ip, &result)
		if err != nil || result == nil {
			return nil, ErrNoDataPresent
		}
		return json.Marshal(result)
	}
	return nil, ErrNoDataPresent
}
