package iqlutil_test

import (
	"testing"

	"github.com/stackql/stackql/internal/stackql/iqlutil"
)

func TestTranslateLikeToRegexPattern(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "no special chars", in: "abc", want: "^abc$"},
		{name: "single wildcard", in: "abc%", want: "^abc.*$"},
		{name: "trailing wildcard", in: "%abc", want: "^.*abc$"},
		{name: "both wildcards", in: "%abc%", want: "^.*abc.*$"},
		{name: "multiple wildcards", in: "a%b%c", want: "^a.*b.*c$"},
		{name: "empty string", in: "", want: "^$"},
		{name: "escaped special regex char", in: "a.b", want: "^a\\.b$"},
		{name: "escaped plus", in: "a+b", want: "^a\\+b$"},
		{name: "escaped dot and wildcard", in: "a.b%", want: "^a\\.b.*$"},
		{name: "only wildcard", in: "%", want: "^.*$"},
		{name: "literal percent in middle", in: "a%b", want: "^a.*b$"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := iqlutil.TranslateLikeToRegexPattern(tc.in)
			if got != tc.want {
				t.Errorf("TranslateLikeToRegexPattern(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSanitisePossibleTickEscapedTerm(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "no ticks", in: "foo", want: "foo"},
		{name: "leading tick", in: "`foo", want: "foo"},
		{name: "trailing tick", in: "foo`", want: "foo"},
		{name: "both ticks", in: "`foo`", want: "foo"},
		{name: "empty string", in: "", want: ""},
		{name: "only leading tick", in: "`", want: ""},
		{name: "only trailing tick", in: "`", want: ""},
		{name: "tick in middle", in: "fo`o", want: "fo`o"},
		{name: "multiple ticks at start", in: "``foo", want: "`foo"},
		{name: "multiple ticks at end", in: "foo``", want: "foo`"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := iqlutil.SanitisePossibleTickEscapedTerm(tc.in)
			if got != tc.want {
				t.Errorf("SanitisePossibleTickEscapedTerm(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestPrettyPrintSomeJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty object", in: `{}`, want: "{}"},
		{name: "empty array", in: `[]`, want: "[]"},
		{name: "single key", in: `{"a":1}`, want: "{\n  \"a\": 1\n}"},
		{name: "single element array", in: `[1]`, want: "[\n  1\n]"},
		{name: "nested object", in: `{"a":{"b":2}}`, want: "{\n  \"a\": {\n    \"b\": 2\n  }\n}"},
		{name: "already pretty", in: "{\n  \"a\": 1\n}", want: "{\n  \"a\": 1\n}"},
		{name: "array of strings", in: `["x","y"]`, want: "[\n  \"x\",\n  \"y\"\n]"},
		{name: "compact object", in: `{"a":1,"b":2}`, want: "{\n  \"a\": 1,\n  \"b\": 2\n}"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := iqlutil.PrettyPrintSomeJSON([]byte(tc.in))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("PrettyPrintSomeJSON(%q) = %q, want %q", tc.in, string(got), tc.want)
			}
		})
	}
}

func TestPrettyPrintSomeJSONError(t *testing.T) {
	_, err := iqlutil.PrettyPrintSomeJSON([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestGetSortedKeysStringMap(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]string
		want []string
	}{
		{name: "empty map", in: map[string]string{}, want: []string{}},
		{name: "single key", in: map[string]string{"a": "1"}, want: []string{"a"}},
		{name: "multiple keys sorted", in: map[string]string{"c": "3", "a": "1", "b": "2"}, want: []string{"a", "b", "c"}},
		{name: "already sorted", in: map[string]string{"a": "1", "b": "2", "c": "3"}, want: []string{"a", "b", "c"}},
		{name: "reverse sorted", in: map[string]string{"c": "3", "b": "2", "a": "1"}, want: []string{"a", "b", "c"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := iqlutil.GetSortedKeysStringMap(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("GetSortedKeysStringMap(%v) = %v, want %v", tc.in, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("GetSortedKeysStringMap(%v)[%d] = %q, want %q", tc.in, i, got[i], tc.want[i])
				}
			}
		})
	}
}