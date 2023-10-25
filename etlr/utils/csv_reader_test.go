// Copyright © 2023 OSINTAMI. This is not yours.
package utils

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func csvReader(t *testing.T, delimiter rune) {
	csvFile := "/tmp/csv-reader-test.csv"
	os.Remove(csvFile)

	var line string
	switch delimiter {
	case '\t':
		line = "# comment line\n\" 1.2.3.4\"\ttrue\n  nope \t\t false \n"
	case ';':
		line = "# comment line\n\" 1.2.3.4\"; true\n  nope ; false \n"
	case ',':
		line = "# comment line\n\" 1.2.3.4\",true\n  nope , false \n"
	case ':':
		line = "# comment line\n\" 1.2.3.4\":true\n  nope : false: extra \n"
	}

	os.WriteFile(csvFile, []byte(line), 0644)

	network := NewNetworkingHelper(nil)

	reader := NewCSVReader()
	err := reader.ProcessFile(csvFile, delimiter, func(values []string) error {
		assert.False(t, strings.Contains(values[0], "#"))
		cidr, err := network.ParseCIDR(values[0])

		switch values[0] {
		case "1.2.3.4":
			assert.Nil(t, err)
			assert.NotNil(t, cidr)
			return nil
		case "nope":
			assert.NotNil(t, err)
			assert.Nil(t, cidr)
			return err
		default:
			assert.Fail(t, "unexpected value")
		}
		return errors.New("cidr invallid")
	})

	assert.Nil(t, err)

	os.Remove(csvFile)
}

func TestCsvReaderTabDelimited(t *testing.T) {
	csvReader(t, '\t')
}

func TestCsvReaderCommaDelimited(t *testing.T) {
	csvReader(t, ',')
}

func TestCsvReaderSemiColonDelimited(t *testing.T) {
	csvReader(t, ';')
}

func TestCsvReaderInvalidFile(t *testing.T) {
	reader := NewCSVReader()
	err := reader.ProcessFile("", ' ', func(values []string) error {
		return nil
	})
	assert.NotNil(t, err)
}

func TestCsvReaderMismatchedColumns(t *testing.T) {
	csvFile := "/tmp/csv-reader-test.csv"
	os.Remove(csvFile)
	var line = "# comment line\n\"1.2.3.4\":true\nnope:false:extra \n"
	os.WriteFile(csvFile, []byte(line), 0644)

	reader := NewCSVReader()
	err := reader.ProcessFile(csvFile, ':', func(values []string) error {
		return nil
	})

	assert.NotNil(t, err)

	os.Remove(csvFile)
}
