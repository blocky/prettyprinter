package prettyprinter_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/blocky/prettyprinter"
	"github.com/blocky/prettyprinter/mocks"
)

func TestPrettyPrinter(t *testing.T) {
	var byteType = mock.AnythingOfType("[]uint8")

	t.Run("happy path", func(t *testing.T) {
		var stderrResult, stdoutResult []byte

		stderr := new(mocks.Writer)
		stdout := new(mocks.Writer)
		stderr.
			On("Write", byteType).
			Run(func(args mock.Arguments) {
				stderrResult = args.Get(0).([]byte)
			}).
			Return(0, nil)
		stdout.
			On("Write", byteType).
			Run(func(args mock.Arguments) {
				stdoutResult = args.Get(0).([]byte)
			}).
			Return(0, nil)

		p := prettyprinter.NewPrettyPrinterFromRaw(stderr, stdout)

		expected := "{\n \"im not happy!\": 18446744073709551615," +
			"\n \"temper tantrum!\": {\n  \"temper\": \"muy malo\"\n }\n}\n"
		assert.Nil(t, p.
			Add(struct {
				F1 uint `json:"im not happy!"`
				F2 struct {
					Temper string `json:"temper"`
				} `json:"temper tantrum!"`
			}{
				F1: ^uint(0),
				F2: struct {
					Temper string `json:"temper"`
				}{"muy malo"},
			}).
			StderrDump().Error(),
		)
		assert.Equal(t, expected, string(stderrResult))
		stderr.AssertExpectations(t)

		expected = "{\n \"hello\": \"world!\"\n}\n"
		assert.Nil(t, p.
			Add(struct {
				Hello string `json:"hello"`
			}{"world!"}).
			StdoutDump().Error(),
		)
		assert.Equal(t, expected, string(stdoutResult))
		stdout.AssertExpectations(t)
	})
	t.Run("writer error", func(t *testing.T) {
		stderr := new(mocks.Writer)
		stdout := new(mocks.Writer)
		stderrErr := errors.New("stderr err!")
		stdoutErr := errors.New("stdout err!")
		stderr.On("Write", byteType).Return(0, stderrErr)
		stdout.On("Write", byteType).Return(0, stdoutErr)

		p := prettyprinter.NewPrettyPrinterFromRaw(stderr, stdout)

		assert.Equal(t, p.
			Add(struct{ BadData int }{1234}).
			StdoutDump().Error(),
			stdoutErr,
		)
		stdout.AssertExpectations(t)

		assert.Equal(t, p.
			Add(prettyprinter.MakeKeyValueError(
				errors.New("cant perform 'rm -rf /usr'"),
			)).
			StderrDump().Error(),
			stderrErr,
		)
		stdout.AssertExpectations(t)
	})
	t.Run("stderr dump on err", func(t *testing.T) {
		stderr := new(mocks.Writer)
		stdout := new(mocks.Writer)
		someErr := errors.New("some err...?")
		expected := "{\n \"err\": \"some err...?\"\n}\n"

		stderr.On("Write", []uint8(expected)).Return(0, nil)
		stdout.On("Write", byteType).Return(0, someErr)

		p := prettyprinter.NewPrettyPrinterFromRaw(stderr, stdout)

		assert.Nil(t, p.
			Add(struct{ Something string }{"abcd"}).
			StdoutDump().
			StderrDumpOnError(),
		)
		stdout.AssertExpectations(t)
		stderr.AssertExpectations(t)
	})
}
