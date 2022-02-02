package models_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestGetErrorCode(t *testing.T) {
	t.Parallel()

	if os.Getenv("TEST_PROCESS") != "1" {
		_, _ = fmt.Fprint(os.Stdout, "TEST_PROCESS is not 1")

		t.SkipNow()
	}

	codeStr := os.Args[2]
	code, _ := strconv.Atoi(codeStr)

	_, _ = fmt.Fprintf(os.Stdout, "exit code will be %d (args: %v)", code, os.Args)

	os.Exit(code)
}

// nolint: gosec
func TestGetResultStatus(t *testing.T) {
	type expStruct struct {
		result models.MutationResult
		err    error
	}

	type tc struct {
		name string
		in   error
		exp  expStruct
	}

	cmd := exec.Command(os.Args[0], "run=TestGetErrorCode", "1")
	cmd.Env = []string{"TEST_PROCESS=1"}
	errCode1 := cmd.Run()
	cmd = exec.Command(os.Args[0], "run=TestGetErrorCode", "2")
	cmd.Env = []string{"TEST_PROCESS=1"}
	errCode2 := cmd.Run()
	cmd = exec.Command(os.Args[0], "run=TestGetErrorCode", "234")
	cmd.Env = []string{"TEST_PROCESS=1"}
	errCode234 := cmd.Run()

	errRandomError := errors.New("some error") // nolint: goerr113

	tcs := []tc{
		{
			name: "escaped mutant - 1",
			in:   &exec.ExitError{ProcessState: &os.ProcessState{}},
			exp: expStruct{
				result: models.ResultMutantEscaped,
				err:    nil,
			},
		},
		{
			name: "escaped mutant - 2",
			in:   errRandomError,
			exp: expStruct{
				result: models.ResultMutantEscaped,
				err:    errRandomError,
			},
		},
		{
			name: "killed mutant",
			in:   errCode1,
			exp: expStruct{
				result: models.ResultMutantKilled,
				err:    nil,
			},
		},
		{
			name: "skipped mutant by code 2",
			in:   errCode2,
			exp: expStruct{
				result: models.ResultMutantSkipped,
				err:    errCode2,
			},
		},
		{
			name: "skipped mutant by default",
			in:   errCode234,
			exp: expStruct{
				result: models.ResultMutantSkipped,
				err:    nil,
			},
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := models.GetResultStatus(tc.in)
			if tc.exp.err != nil {
				if !assert.Error(t, err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, tc.exp.result, out)
		})
	}
}
