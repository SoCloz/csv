package csv

import (
	"bytes"
	"database/sql"
	"os"
	"strconv"
	"testing"
)

func Example() {
	v := []struct {
		Name string
		Age  int
	}{
		{Name: "Bob", Age: 42},
		{Name: "Joe", Age: 17},
	}
	err := NewEncoder(os.Stdout).Encode(v)
	if err != nil {
		panic(err)
	}
	// Output:
	// Name,Age
	// Bob,42
	// Joe,17
}

// Int is an even nuller nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int struct {
	sql.NullInt64
}

func NewInt(i int64, valid bool) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

func TestEncodeNullType(t *testing.T) {
	v := []struct {
		A string
		B Int `csv:"Bis"`
		C bool
		d struct{}
	}{
		{A: "a", B: NewInt(1, true), C: true},
		{A: `b"`, B: NewInt(2, true), C: false},
		{A: `c,`, B: NewInt(3, true), C: true},
	}
	want := "A,Bis,C\n" +
		"a,1,true\n" +
		`"b""",2,false` + "\n" +
		`"c,",3,true` + "\n"

	testEncoding(t, v, want)
}

func TestEncodeSliceStructs(t *testing.T) {
	v := []struct {
		A string
		B int `csv:"Bis"`
		C bool
		d struct{}
	}{
		{A: "a", B: 1, C: true},
		{A: `b"`, B: 2, C: false},
		{A: `c,`, B: 3, C: true},
	}
	want := "A,Bis,C\n" +
		"a,1,true\n" +
		`"b""",2,false` + "\n" +
		`"c,",3,true` + "\n"

	testEncoding(t, v, want)
}

func TestEncodeInterface(t *testing.T) {
	v := []interface{}{struct{ A string }{A: "a"}}
	want := "A\na\n"

	testEncoding(t, v, want)
}

func TestEncodeEmptyStruct(t *testing.T) {
	v := []struct {
		A struct{}
		B int
	}{
		{
			struct{}{},
			1,
		},
	}
	want := "A,B\n,1\n"

	testEncoding(t, v, want)
}

func TestOptions(t *testing.T) {
	v := []struct {
		A string
		B int
	}{
		{A: "a", B: 1},
		{A: "b", B: 2},
	}
	want := "a;1\r\nb;2\r\n"

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.Writer.Comma = ';'
	enc.Writer.UseCRLF = true
	enc.SkipHeader = true

	testCustomEncoding(t, v, want, buf, enc)
}

func testEncoding(t *testing.T, v interface{}, want string) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)

	testCustomEncoding(t, v, want, buf, enc)
}

func testCustomEncoding(t *testing.T, v interface{}, want string, buf *bytes.Buffer, enc *Encoder) {
	err := enc.Encode(v)
	if err != nil {
		t.Fatalf("Encode(%v): %s", v, err)
	}

	got := buf.String()
	if got != want {
		t.Errorf("Invalid encoded value\n got:%s\nwant:%s\n", got, want)
	}
}
