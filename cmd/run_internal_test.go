package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	registerFlags()

	os.Exit(m.Run())
}

func TestRun(t *testing.T) {
	type inStruct struct {
		root string
		args []string
	}

	type expStruct struct {
		results string
		err     error
	}

	type tc struct {
		name string
		in   inStruct
		exp  expStruct
	}

	testArgs := []string{"--exec-timeout=1", "--jobs=1"}
	tcs := []tc{
		{
			name: "success run example tests",
			in: inStruct{
				root: "../example",
				args: testArgs,
			},
			exp: expStruct{
				results: "The mutation score is 0.450000 (total is 20 mutants: 9 killed, 11 escaped, 0 skipped, 8 duplicates)",
				err:     nil,
			},
		},
		{
			name: "success run example tests recursive",
			in: inStruct{
				root: "../example/...",
				args: testArgs,
			},
			exp: expStruct{
				results: "The mutation score is 0.476190 (total is 21 mutants: 10 killed, 11 escaped, 0 skipped, 8 duplicates)",
				err:     nil,
			},
		},
		{
			name: "success run example tests from another directory",
			in: inStruct{
				root: "github.com/smgladkovskiy/go-mutesting/example",
				args: testArgs,
			},
			exp: expStruct{
				results: "The mutation score is 0.450000 (total is 20 mutants: 9 killed, 11 escaped, 0 skipped, 8 duplicates)",
				err:     nil,
			},
		},
		{
			name: "success run example tests with match",
			in: inStruct{
				root: "../example/...",
				args: []string{"--exec=../scripts/exec/test-mutated-package.sh", "--exec-timeout=1", "--match=baz", "./..."},
			},
			exp: expStruct{
				results: "The mutation score is 0.500000 (total is 2 mutants: 1 killed, 1 escaped, 0 skipped, 0 duplicates)",
				err:     nil,
			},
		},
		{
			name: "success run example tests with score",
			in: inStruct{
				root: "../example",
				args: []string{"--exec-timeout=1", "--jobs=1", "--score=0.46"},
			},
			exp: expStruct{
				results: "The mutation score is 0.450000 (total is 20 mutants: 9 killed, 11 escaped, 0 skipped, 8 duplicates)",
				err:     errs.ErrTestScoreWasNotReached,
			},
		},
	}

	// t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()

			saveStderr := os.Stderr
			saveStdout := os.Stdout
			saveCwd, err := os.Getwd()
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			r, w, err := os.Pipe()
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			os.Stderr = w
			os.Stdout = w
			// assert.Nil(t, os.Chdir(root))

			bufChannel := make(chan string)

			go func() {
				buf := new(bytes.Buffer)
				_, err = io.Copy(buf, r)
				assert.NoError(t, err)
				assert.NoError(t, r.Close())

				bufChannel <- buf.String()
			}()

			if err := rootCmd.PersistentFlags().Parse(tc.in.args); !assert.NoError(t, err) {
				t.FailNow()
			}
			rootCmd.SetArgs([]string{tc.in.root})

			cmdErr := rootCmd.Execute()

			assert.Nil(t, w.Close())

			os.Stderr = saveStderr
			os.Stdout = saveStdout
			assert.Nil(t, os.Chdir(saveCwd))

			out := <-bufChannel

			if tc.exp.err != nil {
				assert.ErrorIs(t, cmdErr, tc.exp.err)
			}

			assert.Contains(t, out, tc.exp.results)
		})
	}
}
