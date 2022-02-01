package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
)

func SaveAST(mutationBlackList map[string]struct{}, file string, fset *token.FileSet, node ast.Node) (string, bool, error) {
	var buf bytes.Buffer

	h := md5.New() // nolint: gosec

	if err := printer.Fprint(io.MultiWriter(h, &buf), fset, node); err != nil {
		return "", false, fmt.Errorf("SaveAST dublicate writers error: %w", err)
	}

	checksum := fmt.Sprintf("%x", h.Sum(nil))

	if _, ok := mutationBlackList[checksum]; ok {
		return checksum, true, nil
	}

	mutationBlackList[checksum] = struct{}{}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return "", false, fmt.Errorf("SaveAST source formatting error: %w", err)
	}

	if err := ioutil.WriteFile(file, src, 0o644); err != nil {
		return "", false, fmt.Errorf("SaveAST file writing: %w", err)
	}

	return checksum, false, nil
}
