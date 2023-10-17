// Copyright © 2023 OSINTAMI. This is not yours.
package utils

import (
	"compress/gzip"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/mholt/archiver"
	"github.com/osintami/fingerprintz/log"
)

type IFileSystem interface {
	UnzipFile(zipFile, destination string) error
	UnTarFile(tarFile, destination string) error
	UnGzipFile(gzipPath, gzipFile, outputPath string) error
	Copy(string, string) error
	Move(string, string) error
	Chmod(string, fs.FileMode) error
	Chown(string, string) error
}

type Filesystem struct {
}

func NewFSHelper() *Filesystem {
	return &Filesystem{}
}

func (x *Filesystem) Copy(inFile, outFile string) error {
	src, err := os.Open(inFile)
	if err != nil {
		log.Error().Err(err).Str("component", "filesystem").Str("file", outFile).Msg("open")
		return err
	}
	defer src.Close()

	dst, err := os.Create(outFile)
	if err != nil {
		log.Error().Err(err).Str("component", "filesystem").Str("file", outFile).Msg("create")
		return err
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil || written == 0 {
		log.Error().Err(err).Str("component", "filesystem").Str("file", outFile).Msg("copy")
		if err == nil {
			err = errors.New("no bytes copied")
		}
		return err
	}

	return x.Chmod(outFile, 0644)
}

func (x *Filesystem) Move(inFile, outFile string) error {
	err := os.Rename(inFile, outFile)
	if err != nil {
		log.Error().Err(err).Str("component", "filesystem").Str("file", outFile).Msg("rename")
		return err
	}
	return x.Chmod(outFile, 0644)
}

func (x *Filesystem) Chown(inFile, owner string) error {
	group, err := user.Lookup(owner)
	if err != nil {
		log.Error().Err(err).Str("component", "filesystem").Str("file", inFile).Str("user", owner).Msg("user does not exist")
		return err
	}
	uid, _ := strconv.Atoi(group.Uid)
	gid, _ := strconv.Atoi(group.Gid)

	return syscall.Chown(inFile, uid, gid)
}

func (x *Filesystem) Chmod(inFile string, mode fs.FileMode) error {
	err := os.Chmod(inFile, 0644)
	if err != nil {
		log.Error().Err(err).Str("component", "filesystem").Str("file", inFile).Msg("chmod")
		return err
	}
	return nil
}

func (x *Filesystem) UnzipFile(zipFile, outPath string) error {
	format := archiver.Zip{}
	defer os.Remove(zipFile)
	return format.Unarchive(zipFile, outPath)
}

func (x *Filesystem) UnTarFile(tarFile, outPath string) error {
	format := archiver.Tar{}
	defer os.Remove(tarFile)
	return format.Unarchive(tarFile, outPath)
}

func (x *Filesystem) UnGzipFile(gzipPath, gzipFile, outputPath string) error {

	src, err := os.OpenFile(gzipPath+gzipFile, os.O_RDONLY, 0600)
	if err != nil {
		log.Error().Err(err).Str("component", "filesystem").Str("file", gzipFile).Msg("open input file")
		return err
	}
	defer src.Close()
	gz, err := gzip.NewReader(src)
	if err != nil {
		log.Error().Err(err).Str("component", "filesystem").Str("file", gzipFile).Msg("new reader")
		return err
	}
	defer gz.Close()
	out, err := os.OpenFile(outputPath+strings.TrimRight(gzipFile, ".gz"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error().Err(err).Str("component", "filesystem").Str("file", gzipFile).Msg("open output file")
		return err
	}
	defer out.Close()

	written, err := io.Copy(out, gz)
	if err != nil || written == 0 {
		log.Error().Err(err).Str("component", "filesystem").Str("file", gzipFile).Msg("copy file")
		if err == nil {
			err = errors.New("no bytes copied")
		}
		return err
	}

	os.Remove(gzipPath + gzipFile)

	return nil
}
