package service

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"mash/pkg/unidiff"
	"os"
	"path/filepath"
)

func (svc *Service) walkTarGz(writer *io.PipeWriter, viewPath string, allower *unidiff.Allower) func() error {
	return func() error {
		defer writer.Close()

		gzipWriter := gzip.NewWriter(writer)
		defer gzipWriter.Close()

		tarWriter := tar.NewWriter(gzipWriter)
		defer tarWriter.Close()

		err := filepath.Walk(viewPath,
			func(absPath string, info os.FileInfo, err error) error {
				if err != nil {
					return fmt.Errorf("failed to walk: %w", err)
				}

				if err := addFileToTarWriter(allower, viewPath, absPath, info, tarWriter); err != nil {
					return fmt.Errorf("failed to add file to tar: %w", err)
				}

				return nil
			})
		if err != nil {
			return fmt.Errorf("filewalk failed: %w", err)
		}
		return nil
	}
}

func addFileToTarWriter(allower *unidiff.Allower, viewPath, absPath string, info os.FileInfo, tarWriter *tar.Writer) error {
	if info.IsDir() {
		return nil
	}

	cleanAbsPath := absPath[len(viewPath)+1:]

	if !allower.IsAllowed(cleanAbsPath, info.IsDir()) {
		return nil
	}

	header := &tar.Header{
		Name:    cleanAbsPath,
		Size:    info.Size(),
		Mode:    int64(info.Mode()),
		ModTime: info.ModTime(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(tarWriter, file); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
