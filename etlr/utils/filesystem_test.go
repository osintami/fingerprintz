// Copyright © 2023 OSINTAMI. This is not yours.
package utils

import (
	"io/fs"
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesystemCopy(t *testing.T) {
	testFile := "/tmp/test.csv"

	// cleanup
	os.Remove(testFile)

	// setup
	err := os.WriteFile(testFile, []byte("test"), 0600)
	assert.Nil(t, err)

	outFile := "/tmp/test-copy.csv"
	err = NewFSHelper().Copy(testFile, outFile)
	assert.Nil(t, err)
	assert.True(t, fileExists(outFile))
	assert.True(t, fileExists(testFile))

	// cleanup
	os.Remove(testFile)
	os.Remove(outFile)
}

func TestFilesystemCopyBadInputPath(t *testing.T) {
	err := NewFSHelper().Copy("", "")
	assert.NotNil(t, err)
}

func TestFilesystemCopyBadOutputPath(t *testing.T) {
	testFile := "/tmp/test.csv"

	// cleanup
	os.Remove(testFile)

	// setup
	err := os.WriteFile(testFile, []byte("test"), 0600)
	assert.Nil(t, err)

	err = NewFSHelper().Copy(testFile, "")
	assert.NotNil(t, err)

	// cleanup
	os.Remove(testFile)
}

func TestFilesystemCopyEmptyFile(t *testing.T) {
	testFile := "/tmp/empty.in"
	os.Remove(testFile)

	// setup
	err := os.WriteFile(testFile, make([]byte, 0), 0600)
	assert.Nil(t, err)

	outFile := "/tmp/empty.out"

	err = NewFSHelper().Copy(testFile, outFile)
	assert.NotNil(t, err)

	os.Remove(testFile)
	os.Remove(outFile)
}

func TestFilesystemMove(t *testing.T) {
	testFile := "/tmp/test.csv"
	os.Remove(testFile)
	err := os.WriteFile(testFile, []byte("test"), 0600)
	assert.Nil(t, err)

	outFile := "/tmp/test-copy.csv"
	err = NewFSHelper().Move(testFile, outFile)
	assert.Nil(t, err)
	assert.True(t, fileExists(outFile))
	assert.False(t, fileExists(testFile))

	// bogus input file
	err = NewFSHelper().Move("nope", outFile)
	assert.NotNil(t, err)

	os.Remove(outFile)
}

func TestFilesystemChmod(t *testing.T) {
	testFile := "/tmp/test.csv"
	os.Remove(testFile)
	err := os.WriteFile(testFile, []byte("test"), 0600)
	assert.Nil(t, err)

	info, err := os.Stat(testFile)
	assert.Nil(t, err)
	assert.Equal(t, fs.FileMode(0600), info.Mode().Perm())

	err = NewFSHelper().Chmod(testFile, 0644)
	assert.Nil(t, err)

	info, err = os.Stat(testFile)
	assert.Nil(t, err)
	assert.Equal(t, fs.FileMode(0644), info.Mode().Perm())

	// bogus input file
	err = NewFSHelper().Chmod("nope", 0644)
	assert.NotNil(t, err)

	os.Remove(testFile)
}

func TestFilesystemChown(t *testing.T) {
	testFile := "/tmp/test.csv"
	os.Remove(testFile)
	err := os.WriteFile(testFile, []byte("test"), 0600)
	assert.Nil(t, err)

	uname, err := user.Current()
	assert.Nil(t, err)

	err = NewFSHelper().Chown(testFile, uname.Username)
	assert.Nil(t, err)

	// TODO: check the file's owner

	// bogus input file
	err = NewFSHelper().Chown(testFile, "nope")
	assert.NotNil(t, err)

	os.Remove(testFile)
}

func TestFilesystemUnzipFile(t *testing.T) {
	runDecompress(t, "zip", "test.zip")
}

func TestFilesystemUnTarFile(t *testing.T) {
	runDecompress(t, "tar", "test.tar")
}

func TestFilesystemUnGzipFile(t *testing.T) {
	runDecompress(t, "gz", "test.csv.gz")
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}

func runDecompress(t *testing.T, ext string, testFile string) {
	// setup
	outFile := "/tmp/test.csv"
	os.Remove(outFile)

	fs := NewFSHelper()

	err := fs.Copy("./test/"+testFile, "/tmp/"+testFile)
	assert.Nil(t, err)

	switch ext {
	case "gz":
		err = fs.UnGzipFile("/tmp/", testFile, "/tmp/")
	case "tar":
		err = fs.UnTarFile("/tmp/"+testFile, "/tmp/")
	case "zip":
		err = fs.UnzipFile("/tmp/"+testFile, "/tmp/")
	}
	assert.Nil(t, err)
	assert.True(t, fileExists(outFile))

	err = os.Chmod(outFile, 01777)
	assert.Nil(t, err)

	data, err := os.ReadFile(outFile)
	assert.Nil(t, err)
	assert.Equal(t, "test\n", string(data))

	// teardown
	os.Remove("/tmp/" + testFile)
	os.Remove(outFile)
}

func TestFilesystemUnGzipBadInputPath(t *testing.T) {
	fs := NewFSHelper()
	err := fs.UnGzipFile("/nope/", "", "/tmp/")
	assert.NotNil(t, err)
}

func TestFilesystemUnGzipCorruptGzipFile(t *testing.T) {
	fs := NewFSHelper()
	inFile := "./test/test.csv.gz"
	outFile := "/tmp/test.csv.gz"

	// cleanup
	os.Remove(outFile)

	// setup test
	err := fs.Copy(inFile, outFile)
	assert.Nil(t, err)

	err = fs.UnGzipFile("./test/", "corrupt.gz", "/tmp/")
	assert.NotNil(t, err)

	// cleanup
	os.Remove(outFile)
}

func TestFilesystemUnGzipBadOutputPath(t *testing.T) {
	fs := NewFSHelper()
	inFile := "./test/test.csv.gz"
	outFile := "/tmp/test.csv.gz"

	// cleanup
	os.Remove(outFile)

	// setup test
	err := fs.Copy(inFile, outFile)
	assert.Nil(t, err)

	// bad output path
	err = fs.UnGzipFile("/tmp/", "test.csv.gz", "/nope/")
	assert.NotNil(t, err)

	// cleanup
	os.Remove(outFile)
}

func TestFilesystemUnGzipEmptyFile(t *testing.T) {
	fs := NewFSHelper()
	inFile := "./test/empty.gz"
	outFile := "/tmp/empty.gz"

	// cleanup
	os.Remove(outFile)

	// setup test
	err := fs.Copy(inFile, outFile)
	assert.Nil(t, err)

	// empty archive triggers ****** stream copy fail
	err = fs.UnGzipFile("/tmp/", "empty.gz", "/tmp/")
	assert.NotNil(t, err)

	// cleanup
	os.Remove("/tmp/empty")
	os.Remove(outFile)
}
