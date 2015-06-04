// Package csv provides a CSV Encoder similar to json.Encoder or xml.Encoder.
package csv

import (
	"encoding"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// An Encoder writes CSV objects to an output stream.
type Encoder struct {
	SkipHeader bool
	Writer     *csv.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	enc := &Encoder{Writer: csv.NewWriter(w)}

	return enc
}

// Encode writes the CSV encoding of v to the stream.
func (e *Encoder) Encode(v interface{}) error {
	val := reflect.ValueOf(v)
	t := val.Type()

	if k := t.Kind(); k != reflect.Slice {
		return errors.New("csv: unsupported type " + k.String())
	}

	var indexes []int
	for i := 0; i < val.Len(); i++ {
		s := getPointee(val.Index(i))
		if k := s.Type().Kind(); k != reflect.Struct {
			return errors.New("csv: unsupported type slice of " +
				k.String() + "s")
		}
		if i == 0 {
			indexes = getExpFieldIndexes(s.Type())

			if !e.SkipHeader {
				header := getStructHeader(s.Type(), indexes)
				if err := e.Writer.Write(header); err != nil {
					return err
				}
			}
		}
		line := structToStrings(s, indexes)

		if err := e.Writer.Write(line); err != nil {
			return err
		}
	}

	e.Writer.Flush()

	return nil
}

func getPointee(v reflect.Value) reflect.Value {
	for {
		switch v.Type().Kind() {
		case reflect.Interface, reflect.Ptr:
			v = v.Elem()
		default:
			return v
		}
	}
}

func getExpFieldIndexes(t reflect.Type) []int {
	n := 0
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath == "" {
			n++
		}
	}

	indexes := make([]int, n)
	i := 0
	for id := 0; id < t.NumField(); id++ {
		if t.Field(id).PkgPath == "" {
			indexes[i] = id
			i++
		}
	}

	return indexes
}

func getStructHeader(t reflect.Type, indexes []int) []string {
	h := make([]string, len(indexes))
	for _, i := range indexes {
		f := t.Field(i)
		if tag := f.Tag.Get("csv"); tag != "" {
			h[i] = tag
		} else {
			h[i] = f.Name
		}
	}

	return h
}

var (
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func structToStrings(s reflect.Value, indexes []int) []string {
	rows := make([]string, len(indexes))
	for _, i := range indexes {
		f := s.Field(i)
		switch f.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice:
			if !f.IsNil() {
				rows[i] = fmt.Sprintf("%v", f.Interface())
			}
		case reflect.Struct:
			v := f.Interface()
			if v != struct{}{} {
				if f.CanInterface() && f.Type().Implements(textMarshalerType) {
					text, _ := f.Interface().(encoding.TextMarshaler).MarshalText()
					rows[i] = string(text)
				} else {
					rows[i] = fmt.Sprintf("%v", v)
				}
			}
		default:
			rows[i] = fmt.Sprintf("%v", f.Interface())
		}
	}

	return rows
}
