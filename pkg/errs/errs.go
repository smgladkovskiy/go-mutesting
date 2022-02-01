package errs

import (
	"errors"
)

var (
	ErrNoSuitableSourceFiles  = errors.New("could not find any suitable Go source files")
	ErrNotMD5Checksum         = errors.New("not a MD5 checksum")
	ErrTestScoreWasNotReached = errors.New("test score was not reached")
	ErrNoExecWasExecuted      = errors.New("cannot do a mutation testing summary since no exec command was executed")

	ErrUnknownMutator = errors.New("unknown mutator")

	ErrWrongExitCode = errors.New("wrong exit code")

	ErrNilMutatorFunction = errors.New("mutator function is nil")
	ErrMutatorRegistered  = errors.New("mutator function is nil")
)
