package models

import (
	"errors"
	"os/exec"
	"syscall"
)

type MutationResult int

const (
	ResultMutantEscaped MutationResult = iota
	ResultMutantKilled
	ResultMutantSkipped
	ResultEmpty
)

func GetResultStatus(err error) (MutationResult, error) {
	var e *exec.ExitError

	if errors.As(err, &e) {
		switch e.Sys().(syscall.WaitStatus).ExitStatus() {
		case 0:
			return ResultMutantEscaped, nil
		case 1:
			return ResultMutantKilled, nil
		case 2: // nolint: gomnd
			return ResultMutantSkipped, err
		default:
			return ResultMutantSkipped, nil
		}
	}

	return ResultMutantEscaped, err
}
