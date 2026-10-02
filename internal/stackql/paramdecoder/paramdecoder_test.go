package paramdecoder_test

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"testing"
	"time"

	"github.com/lib/pq/oid"
	"github.com/stackql/stackql/internal/stackql/paramdecoder"
)

func TestNewDecoder(t *testing.T) {
	d := paramdecoder.NewDecoder()
	if d == nil {
		t.Fatal("NewDecoder() returned nil")
	}
	_, ok := d.(paramdecoder.Decoder)
	if !ok {
		t.Fatalf("NewDecoder() returned %T, want Decoder", d)
	}
}

func TestDecodeParamsTextFormat(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Test with all text format (empty format slice)
	cases := []struct {
		name       string
		paramOIDs  []uint32
		paramForms []int16
		paramVals  [][]byte
		want       []string
	}{
		{
			name:     "empty",
			paramOIDs:  []uint32{},
			paramForms: []int16{},
			paramVals:  [][]byte{},
			want:       []string{},
		},
		{
			name:     "single text",
			paramOIDs:  []uint32{uint32(oid.T_text)},
			paramForms: []int16{0}, // text format
			paramVals:  [][]byte{[]byte("hello")},
			want:       []string{"hello"},
		},
		{
			name:     "multiple text",
			paramOIDs:  []uint32{uint32(oid.T_text), uint32(oid.T_text)},
			paramForms: []int16{0, 0},
			paramVals:  [][]byte{[]byte("a"), []byte("b")},
			want:       []string{"a", "b"},
		},
		{
			name:     "with null",
			paramOIDs:  []uint32{uint32(oid.T_text), uint32(oid.T_text)},
			paramForms: []int16{0, 0},
			paramVals:  [][]byte{[]byte("hello"), nil},
			want:       []string{"hello", "NULL"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := d.DecodeParams(tc.paramOIDs, tc.paramForms, tc.paramVals)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("got %d values, want %d", len(got), len(tc.want))
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestDecodeParamBool(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Test bool true (0x01)
	val := []byte{0x01}
	got, err := d.DecodeParams([]uint32{uint32(oid.T_bool)}, []int16{1}, [][]byte{val})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "true" {
		t.Errorf("bool true: got %q, want true", got[0])
	}

	// Test bool false (0x00)
	val = []byte{0x00}
	got, err = d.DecodeParams([]uint32{uint32(oid.T_bool)}, []int16{1}, [][]byte{val})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "false" {
		t.Errorf("bool false: got %q, want false", got[0])
	}

	// Test wrong length
	val = []byte{0x00, 0x00} // 2 bytes instead of 1
	_, err = d.DecodeParams([]uint32{uint32(oid.T_bool)}, []int16{1}, [][]byte{val})
	if err == nil {
		t.Error("expected error for wrong bool length")
	}
}

func TestDecodeParamInt2(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Test int2: 1000 (0x03 0xE8)
	val := []byte{0x03, 0xE8}
	got, err := d.DecodeParams([]uint32{uint32(oid.T_int2)}, []int16{1}, [][]byte{val})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "1000" {
		t.Errorf("int2 1000: got %q, want 1000", got[0])
	}

	// Test int2: -1 (0xFF 0xFF)
	val = []byte{0xFF, 0xFF}
	got, err = d.DecodeParams([]uint32{uint32(oid.T_int2)}, []int16{1}, [][]byte{val})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "-1" {
		t.Errorf("int2 -1: got %q, want -1", got[0])
	}

	// Test wrong length
	val = []byte{0x00, 0x00, 0x00} // 3 bytes instead of 2
	_, err = d.DecodeParams([]uint32{uint32(oid.T_int2)}, []int16{1}, [][]byte{val})
	if err == nil {
		t.Error("expected error for wrong int2 length")
	}
}

func TestDecodeParamInt4(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Test int4: 50000 (0x00 0x00 0xC3 0x50)
	val := []byte{0x00, 0x00, 0xC3, 0x50}
	got, err := d.DecodeParams([]uint32{uint32(oid.T_int4)}, []int16{1}, [][]byte{val})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "50000" {
		t.Errorf("int4 50000: got %q, want 50000", got[0])
	}

	// Test int4: -1 (0xFF 0xFF 0xFF 0xFF)
	val = []byte{0xFF, 0xFF, 0xFF, 0xFF}
	got, err = d.DecodeParams([]uint32{uint32(oid.T_int4)}, []int16{1}, [][]byte{val})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "-1" {
		t.Errorf("int4 -1: got %q, want -1", got[0])
	}
}

func TestDecodeParamFloat4(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Test float4: 3.14 (IEEE 754: 0x40 0x48 0xF5 0xC3)
	val := []byte{0x40, 0x48, 0xF5, 0xC3}
	got, err := d.DecodeParams([]uint32{uint32(oid.T_float4)}, []int16{1}, [][]byte{val})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := strconv.FormatFloat(3.14, 'f', -1, 32)
	if got[0] != expected {
		t.Errorf("float4 3.14: got %q, want %q", got[0], expected)
	}
}

func TestDecodeParamFloat8(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Test float8: 2.718281828459045 (IEEE 754: 0x40 0x05 0xBF 0x0A 0x8B 0x14 0x57 0x69)
	val := []byte{0x40, 0x05, 0xBF, 0x0A, 0x8B, 0x14, 0x57, 0x69}
	got, err := d.DecodeParams([]uint32{uint32(oid.T_float8)}, []int16{1}, [][]byte{val})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := strconv.FormatFloat(2.718281828459045, 'f', -1, 64)
	if got[0] != expected {
		t.Errorf("float8 2.718281828459045: got %q, want %q", got[0], expected)
	}
}

func TestDecodeParamTimestamp(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Test timestamp: 2006-01-02 15:04:05 (microseconds since 2000-01-01)
	pgEpoch := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	target := time.Date(2006, 1, 2, 15, 04, 05, 0, time.UTC)
	microseconds := target.Sub(pgEpoch).Microseconds()
	val := make([]byte, 8)
	binary.BigEndian.PutUint64(val, uint64(microseconds))

	got, err := d.DecodeParams([]uint32{uint32(oid.T_timestamp)}, []int16{1}, [][]byte{val})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "2006-01-02 15:04:05" {
		t.Errorf("timestamp: got %q, want 2006-01-02 15:04:05", got[0])
	}
}

func TestDecodeParamTextTypes(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Test text, varchar, name - all treated as text
	tests := []struct {
		name  string
		oid   uint32
		input string
	}{
		{name: "text", oid: uint32(oid.T_text), input: "hello world"},
		{name: "varchar", oid: uint32(oid.T_varchar), input: "hello world"},
		{name: "name", oid: uint32(oid.T_name), input: "hello world"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := d.DecodeParams([]uint32{tc.oid}, []int16{0}, [][]byte{[]byte(tc.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got[0] != tc.input {
				t.Errorf("%s: got %q, want %q", tc.name, got[0], tc.input)
			}
		})
	}
}

func TestDecodeParamsMixedFormats(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Mixed formats: first text (format 0), second binary (format 1)
	// Binary bool true
	val := []byte{0x01}
	got, err := d.DecodeParams(
		[]uint32{uint32(oid.T_text), uint32(oid.T_bool)},
		[]int16{0, 1}, // text, binary
		[][]byte{[]byte("hello"), val},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "hello" {
		t.Errorf("text param: got %q, want hello", got[0])
	}
	if got[1] != "true" {
		t.Errorf("bool param: got %q, want true", got[1])
	}

	// Single format applies to all - binary format with invalid data for bool
	got, err = d.DecodeParams(
		[]uint32{uint32(oid.T_text), uint32(oid.T_bool)},
		[]int16{1}, // binary format for all
		[][]byte{[]byte("world"), []byte{0x00}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With binary format, "world" as text is treated as text (falls through to string(val))
	// and bool 0x00 is "false"
	if got[0] != "world" {
		t.Errorf("text param in binary format: got %q, want world", got[0])
	}
	if got[1] != "false" {
		t.Errorf("bool param in binary format: got %q, want false", got[1])
	}
}

func TestDecodeParamsErrorCases(t *testing.T) {
	d := paramdecoder.NewDecoder()

	// Wrong int2 length
	_, err := d.DecodeParams([]uint32{uint32(oid.T_int2)}, []int16{1}, [][]byte{[]byte{0x00}})
	if err == nil {
		t.Error("expected error for short int2")
	}

	// Wrong int4 length
	_, err = d.DecodeParams([]uint32{uint32(oid.T_int4)}, []int16{1}, [][]byte{[]byte{0x00, 0x00}})
	if err == nil {
		t.Error("expected error for short int4")
	}

	// Wrong float4 length
	_, err = d.DecodeParams([]uint32{uint32(oid.T_float4)}, []int16{1}, [][]byte{[]byte{0x00, 0x00}})
	if err == nil {
		t.Error("expected error for short float4")
	}

	// Wrong float8 length
	_, err = d.DecodeParams([]uint32{uint32(oid.T_float8)}, []int16{1}, [][]byte{[]byte{0x00}})
	if err == nil {
		t.Error("expected error for short float8")
	}

	// Wrong timestamp length
	_, err = d.DecodeParams([]uint32{uint32(oid.T_timestamp)}, []int16{1}, [][]byte{[]byte{0x00, 0x00, 0x00}})
	if err == nil {
		t.Error("expected error for short timestamp")
	}

	// Verify error message has parameter prefix
	_, err = d.DecodeParams([]uint32{uint32(oid.T_int2)}, []int16{1}, [][]byte{[]byte{0x00}})
	if err != nil && !bytes.Contains([]byte(err.Error()), []byte("parameter")) {
		t.Errorf("error %q does not contain 'parameter'", err.Error())
	}
}
