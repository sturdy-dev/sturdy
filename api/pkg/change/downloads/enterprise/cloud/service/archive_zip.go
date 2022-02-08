package service

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"getsturdy.com/api/pkg/unidiff"
)

func (svc *Service) walkZip(writer *io.PipeWriter, viewPath string, allower *unidiff.Allower) func() error {
	return func() error {
		defer writer.Close()

		zipWriter := zip.NewWriter(writer)
		defer zipWriter.Close()

		err := filepath.Walk(viewPath,
			func(absPath string, info os.FileInfo, err error) error {
				if err != nil {
					return fmt.Errorf("failed to walk: %w", err)
				}

				if err := addFileToZipWriter(allower, viewPath, absPath, info, zipWriter); err != nil {
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

func addFileToZipWriter(allower *unidiff.Allower, viewPath, absPath string, info os.FileInfo, zipWriter *zip.Writer) error {
	if info.IsDir() {
		return nil
	}

	cleanAbsPath := absPath[len(viewPath)+1:]

	if !allower.IsAllowed(cleanAbsPath, info.IsDir()) {
		return nil
	}

	header := &zip.FileHeader{
		Name:     cleanAbsPath,
		Modified: info.ModTime(),
	}

	fileWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(fileWriter, file); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
