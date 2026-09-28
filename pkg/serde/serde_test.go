package serde_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stackql/stackql/pkg/serde"
)

func TestNewStringArrayMapSerDe(t *testing.T) {
	s := serde.NewStringArrayMapSerDe()
	assert.NotNil(t, s)
}

func TestSerialize(t *testing.T) {
	cases := []struct {
		name string
		arr  []string
		want string
	}{
		{name: "empty", arr: []string{}, want: ""},
		{name: "single", arr: []string{"a"}, want: "a"},
		{name: "multiple", arr: []string{"a", "b", "c"}, want: "a,b,c"},
		{name: "empty strings", arr: []string{"", "", ""}, want: ",,"},
	}
	s := serde.NewStringArrayMapSerDe()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Serialize(tc.arr)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDeserialize(t *testing.T) {
	cases := []struct {
		name string
		str  string
		want map[string]any
	}{
		{name: "empty", str: "", want: map[string]any{}},
		{name: "single", str: "a", want: map[string]any{"a": struct{}{}}},
		{name: "multiple", str: "a,b,c", want: map[string]any{"a": struct{}{}, "b": struct{}{}, "c": struct{}{}}},
		{name: "with-empty-elems", str: "a,,b", want: map[string]any{"a": struct{}{}, "b": struct{}{}}},
	}
	s := serde.NewStringArrayMapSerDe()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Deserialize(tc.str)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSerializeDeserializeRoundTrip(t *testing.T) {
	s := serde.NewStringArrayMapSerDe()
	original := []string{"alpha", "beta", "gamma"}
	serialized, err := s.Serialize(original)
	assert.NoError(t, err)
	deserialized, err := s.Deserialize(serialized)
	assert.NoError(t, err)
	for _, item := range original {
		assert.Contains(t, deserialized, item)
	}
}

func TestDeserializeDuplicateKeys(t *testing.T) {
	s := serde.NewStringArrayMapSerDe()
	_, err := s.Deserialize("a,a,b")
	assert.NoError(t, err)
	// map deduplication: "a" appears once
	got, err := s.Deserialize("a,a,b")
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{"a": struct{}{}, "b": struct{}{}}, got)
}
