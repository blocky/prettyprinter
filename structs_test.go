package prettyprinter_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blocky/prettyprinter"
)

type FieldedStruct struct {
	FieldError prettyprinter.FieldError `json:"err"`
}

func makeFieldedStruct(fe prettyprinter.FieldError) FieldedStruct {
	return FieldedStruct{fe}
}

func TestFieldError(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		msg := "yay for formatting errors"
		e := errors.New(msg)
		fe := prettyprinter.FieldError{e}
		fs := makeFieldedStruct(fe)
		expected := fmt.Sprintf("{\n \"err\": \"%s\"\n}", msg)

		b, err := json.MarshalIndent(fs, "", " ")
		assert.Nil(t, err)
		assert.Equal(t, expected, string(b))
	})

	t.Run("given nil error", func(t *testing.T) {
		var e error = nil
		fe := prettyprinter.FieldError{e}
		fs := makeFieldedStruct(fe)
		expected := "{\n \"err\": null\n}"

		b, err := json.MarshalIndent(fs, "", " ")
		assert.Nil(t, err)
		assert.Equal(t, expected, string(b))
	})
}

func TestKeyValueError(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		msg := "nested error value"
		e := errors.New(msg)
		kve := prettyprinter.MakeKeyValueError(e)
		expected := fmt.Sprintf("{\n \"err\": \"%s\"\n}", msg)

		b, err := kve.MarshalJSON()
		assert.Nil(t, err)
		assert.Equal(t, expected, string(b))
	})

	t.Run("given nil error", func(t *testing.T) {
		var e error = nil
		kve := prettyprinter.MakeKeyValueError(e)
		expected := "{\n \"err\": null\n}"

		b, err := kve.MarshalJSON()
		assert.Nil(t, err)
		assert.Equal(t, expected, string(b))
	})
}
