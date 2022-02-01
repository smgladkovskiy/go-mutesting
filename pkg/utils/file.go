package utils

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies a file from src to dst.
func CopyFile(src string, dst string) (err error) {
	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("copyFile os.Open error: %w", err)
	}

	defer func() {
		e := s.Close()
		if err == nil {
			err = e
		}
	}()

	d, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("copyFile os.Create error: %w", err)
	}

	defer func() {
		e := d.Close()
		if err == nil {
			err = e
		}
	}()

	if _, err = io.Copy(d, s); err != nil {
		return fmt.Errorf("copyFile io.Copy error: %w", err)
	}

	i, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("copyFile os.Stat error: %w", err)
	}

	return os.Chmod(dst, i.Mode())
}
