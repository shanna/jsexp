package jsexp

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

/*
  TODO: To solve the symbol/string thing I think passing a predefined list of symbols is the best
  solution. Or perhaps inspect each string and everything that can be a symbol is?
*/

type Encoder struct {
	w    io.Writer
	data interface{}
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) encodeMap(m map[string]interface{}) string {
	values := []string{}
	for k, v := range m {
		values = append(values, fmt.Sprintf("(%s %s)", k, enc.encode(v)))
	}
	return fmt.Sprintf("(%s)", strings.Join(values, " "))
}

func (enc *Encoder) encodeSlice(s []interface{}) string {
	values := []string{}
	for i, v := range s {
		if i == 0 {
			// TODO: Treating the fist value as symbol and don't quote it is a hack job. A better
			// solution is a whitelist of symbols and/or to inspect each string to see if I can
			// treat it as a symbol.
			values = append(values, fmt.Sprintf("%v", v))
		} else {
			values = append(values, enc.encode(v))
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(values, " "))
}

func (enc *Encoder) encode(data interface{}) string {
	switch t := reflect.ValueOf(data); t.Kind() {
	default:
		return fmt.Sprintf("%v", data) // TODO: Explicit types. It'd be nice if it didn't mess with Floats.
	case reflect.String:
		return strconv.Quote(data.(string))
	case reflect.Map:
		return enc.encodeMap(data.(map[string]interface{}))
	case reflect.Slice:
		return enc.encodeSlice(data.([]interface{}))
	}
	return ""
}

// TODO: func (enc *Encoder) Encode(stream []byte) error {

func (enc *Encoder) EncodeJSON(dec *json.Decoder) error {
	for {
		if err := dec.Decode(&enc.data); err == io.EOF {
			break
		} else if err != nil {
			// log.Fatal(err)
		}
		io.WriteString(enc.w, enc.encode(enc.data))
	}
	return nil
}
