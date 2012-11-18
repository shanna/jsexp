package jsexp

import (
	"bytes"
	"encoding/json"
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func convert(s string) (string, error) {
	buf := new(bytes.Buffer)
	dec := json.NewDecoder(strings.NewReader(s))
	enc := NewEncoder(buf)
	err := enc.EncodeJSON(dec)
	return buf.String(), err
}

func assert_equal(t *testing.T, in string, expects string) {
	out, err := convert(in)
	assert.Equal(t, nil, err)
	assert.Equal(t, expects, out)
}

func TestJsexpEncodeString(t *testing.T) {
	assert_equal(t, `"hello"`, `"hello"`)
	assert_equal(t, `"he\"llo"`, `"he\"llo"`)
}

func TestJsexpEncodeNumbers(t *testing.T) {
	assert_equal(t, `0.4`, `0.4`)
	assert_equal(t, `2.4`, `2.4`)

	// The current implementation converts types, I'm not converting a tokenized JSON stream straight into an Sexp stream.
	assert_equal(t, `2.0`, `2`)
	assert_equal(t, `2.0e2`, `200`)
}

func TestJsexpEncodeArray(t *testing.T) {
	assert_equal(t, `["hello", 2]`, `(hello 2)`)
	assert_equal(t, `["hello", "foo"]`, `(hello "foo")`)
	assert_equal(t, `["hello", "foo", "bar"]`, `(hello "foo" "bar")`)
	assert_equal(t, `["foo", ["hello", 2]]`, `(foo (hello 2))`)
}

func TestJsexpEncodeHash(t *testing.T) {
	assert_equal(t, `{"foo": "bar"}`, `((foo "bar"))`)
}
